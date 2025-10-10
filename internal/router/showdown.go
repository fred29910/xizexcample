package router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"../logic"
	"../server"
	"../msg"
	"encoding/json"
)

// ShowdownHandler 处理摊牌请求
type ShowdownHandler struct {
	BaseRouter
}

// Handle 处理请求
func (h *ShowdownHandler) Handle(request ziface.IRequest) {
	// 1. 获取玩家和房间 (这里简化处理)
	playerID := request.GetConnection().GetConnID()
	roomManager := server.GetRoomManager()

	var playerRoom *logic.Room
	var targetPlayer *logic.Player
	for _, room := range roomManager.GetAllRooms() {
		p, err := room.GetPlayer(playerID)
		if err == nil {
			playerRoom = room
			targetPlayer = p
			break
		}
	}

	if targetPlayer == nil {
		h.sendErrorResponse(request.GetConnection(), "Player not found in any room")
		return
	}

	// 2. 检查游戏状态是否允许摊牌
	if playerRoom.GetFSM().GetCurrentState() != logic.STATE_SHOWDOWN {
		h.sendErrorResponse(request.GetConnection(), "Cannot request showdown at this time")
		return
	}

	// 3. 检查是否是庄家请求摊牌
	// 在斗牛牛中，通常由庄家或所有玩家都确认后才能摊牌
	if !targetPlayer.IsBanker() {
		// TODO: 可以添加逻辑，例如需要所有玩家都确认才能摊牌
		// 这里简化处理，允许任何玩家请求
	}

	// 4. 执行摊牌逻辑
	err := playerRoom.GetFSM().Showdown()
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), "Failed to transition to settlement state")
		return
	}

	// 5. 计算所有玩家的牌型
	results := make([]*msg.PlayerResult, 0)
	hands := make([][]*logic.Card, 0)
	for _, player := range playerRoom.GetPlayers() {
		if player.GetStatus() == logic.STATUS_PLAYING {
			hand := player.GetHand()
			hands = append(hands, hand)

			cards := make([]*msg.Card, len(hand))
			for i, card := range hand {
				cards[i] = &msg.Card{
					Suit: uint32(card.Suit),
					Rank: uint32(card.Rank),
				}
			}

			cardType, bullValue := logic.CalculateBull(hand)

			result := &msg.PlayerResult{
				PlayerId: uint64(player.ID),
				Cards:    cards,
				CardType: msg.CardType(cardType),
				BullValue: uint32(bullValue),
				IsWinner: false,
			}
			results = append(results, result)
		}
	}

	// 6. 确定赢家
	// 比较所有玩家的牌，找出赢家
	winnerIndex := -1
	for i := 0; i < len(results); i++ {
		if results[i].IsWinner { // 初始化时，所有玩家都不是赢家
			continue
		}
		isWinner := true
		for j := 0; j < len(results); j++ {
			if i != j {
				// 比较手牌 i 和 j
				// 注意：这里需要将 []*msg.Card 转换回 []*logic.Card
				// 为了简化，我们假设 CompareHands 可以直接比较 msg.Card
				// 或者我们重构一下，直接比较 logic.Card
				handI := hands[i]
				handJ := hands[j]
				if logic.CompareHands(handI, handJ) < 0 { // 如果 handI 比 handJ 小，则不是赢家
					isWinner = false
					break
				}
			}
		}
		if isWinner {
			results[i].IsWinner = true
			winnerIndex = i
			break // 假设只有一个赢家，平局情况暂不处理
		}
	}

	// 7. 广播摊牌结果
	showdownNtf := &msg.S2C_ShowdownNtf{
		Results: results,
	}
	ntfData, err := json.Marshal(showdownNtf)
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), "Failed to marshal showdown notification")
		return
	}

	// TODO: 广播给所有玩家
	// for _, player := range playerRoom.GetPlayers() {
	//     if conn := getPlayerConnection(player.ID); conn != nil {
	//         conn.SendMsg(uint32(msg.S2C_SHOWDOWN_NTF), ntfData)
	//     }
	// }

	// 8. 触发结算
	go h.triggerSettlement(playerRoom, results)

	// 9. 发送确认响应
	showdownAck := &msg.S2C_ShowdownAck{
		RetCode: 0,
	}
	ackData, err := json.Marshal(showdownAck)
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), "Failed to marshal response")
		return
	}
	request.GetConnection().SendMsg(uint32(msg.S2C_SHOWDOWN_ACK), ackData)
}

// triggerSettlement 触发结算逻辑
func (h *ShowdownHandler) triggerSettlement(room *logic.Room, results []*msg.PlayerResult) {
	// 1. 转换到结算状态
	err := room.GetFSM().Settlement()
	if err != nil {
		// TODO: 记录错误
		return
	}

	// 2. 更新玩家分数
	for _, result := range results {
		player, err := room.GetPlayer(int32(result.PlayerId))
		if err != nil {
			continue
		}

		// TODO: 实现真正的分数计算逻辑
		var scoreChange int64
		if result.IsWinner {
			// 赢家获得下注金额
			scoreChange = int64(player.GetBetAmount())
		} else {
			// 输家失去下注金额
			scoreChange = -int64(player.GetBetAmount())
		}
		player.AddScore(scoreChange)
		player.SetBetAmount(0) // 重置下注
	}

	// 3. 广播游戏结果
	gameResultNtf := &msg.S2C_GameResultNtf{
		RoundId: uint64(room.ID), // 使用房间ID作为局次ID
		Results: results,
	}
	ntfData, _ := json.Marshal(gameResultNtf)

	// TODO: 广播给所有玩家
	// for _, player := range room.GetPlayers() {
	//     if conn := getPlayerConnection(player.ID); conn != nil {
	//         conn.SendMsg(uint32(msg.S2C_GAME_RESULT_NTF), ntfData)
	//     }
	// }

	// 4. 广播房间状态更新
	broadcastRoomState(room)
}

// sendErrorResponse 发送错误响应
func (h *ShowdownHandler) sendErrorResponse(conn ziface.IConnection, errorMsg string) {
	errAck := map[string]interface{}{
		"ret_code": 1,
		"message":  errorMsg,
	}
	ackData, _ := json.Marshal(errAck)
	conn.SendMsg(uint32(msg.S2C_SHOWDOWN_ACK), ackData)
}