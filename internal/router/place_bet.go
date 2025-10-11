package router

import (
	"encoding/json"
	"github.com/aceld/zinx/ziface"
	"xizexcample/internal/logic"
	"xizexcample/internal/msg"
	"xizexcample/internal/pkg/logger"
)

// PlaceBetHandler 处理下注请求
type PlaceBetHandler struct {
	BaseRouter
}

// Handle 处理请求
func (h *PlaceBetHandler) Handle(request ziface.IRequest) {
	// 1. 解析客户端请求
	var betReq msg.C2S_PlaceBetReq
	err := json.Unmarshal(request.GetData(), &betReq)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal place bet request: %v", err)
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_PLACE_BET_ACK), "Invalid request data")
		return
	}

	// 2. 获取玩家和房间
	targetPlayer, playerRoom, err := GetPlayerAndRoom(request)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get player and room: %v", err)
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_PLACE_BET_ACK), err.Error())
		return
	}

	// 3. 检查游戏状态是否允许下注
	if playerRoom.GetFSM().GetCurrentState() != logic.STATE_BETTING {
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_PLACE_BET_ACK), "Cannot place bet at this time")
		return
	}

	// 4. 检查玩家是否已经下注
	if targetPlayer.HasBet() {
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_PLACE_BET_ACK), "Player has already placed a bet")
		return
	}

	// 5. 检查下注金额是否有效
	if betReq.Multiple <= 0 {
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_PLACE_BET_ACK), "Invalid bet amount")
		return
	}

	// TODO: 检查玩家余额是否足够
	// if targetPlayer.GetScore() < betReq.Multiple {
	//     sendErrorResponse(request.GetConnection(), "Insufficient balance")
	//     return
	// }

	// 6. 处理下注逻辑
	targetPlayer.PlaceBet(betReq.Multiple)
	logger.InfoLogger.Printf("Player %d in room %d placed a bet of %d", targetPlayer.ID, playerRoom.ID, betReq.Multiple)

	// 7. 广播下注信息
	// broadcastBetInfo(playerRoom, targetPlayer, betReq.Multiple) // TODO: S2C_BetNtf has different fields

	// 8. 检查是否所有玩家都已下注
	allBetsPlaced := true
	for _, p := range playerRoom.GetPlayers() {
		if !p.HasBet() && p.GetStatus() == logic.STATUS_PLAYING {
			allBetsPlaced = false
			break
		}
	}

	// 9. 如果所有玩家都已下注，转换游戏状态到摊牌
	if allBetsPlaced {
		logger.InfoLogger.Printf("All players in room %d have placed their bets", playerRoom.ID)
		err = playerRoom.GetFSM().PlaceBet()
		if err != nil {
			logger.ErrorLogger.Printf("Failed to transition to showdown state in room %d: %v", playerRoom.ID, err)
			sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_PLACE_BET_ACK), "Failed to transition to showdown state")
			return
		}

		// TODO: 自动触发摊牌逻辑
		// go triggerShowdown(playerRoom)
	}

	// 10. 发送确认响应
	betAck := &msg.S2C_PlaceBetAck{
		RetCode:  0,
		Multiple: betReq.Multiple,
	}
	ackData, err := json.Marshal(betAck)
	if err != nil {
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_PLACE_BET_ACK), "Failed to marshal response")
		return
	}
	request.GetConnection().SendMsg(uint32(msg.MsgID_S2C_PLACE_BET_ACK), ackData)

	// 11. 广播房间状态更新
	broadcastRoomState(playerRoom)
}

// broadcastBetInfo 广播下注信息给所有玩家
// func broadcastBetInfo(room *logic.Room, player *logic.Player, amount int32) {
// 	betInfo := &msg.S2C_BetNtf{
// 		BankerId: room.GetBankerID(),
// 		Countdown: 10, // Example
// 	}
// 	ackData, _ := json.Marshal(betInfo)
//
// 	// TODO: 遍历房间内所有玩家的连接并发送消息
// 	// for _, p := range room.GetPlayers() {
// 	//     if conn := getPlayerConnection(p.ID); conn != nil {
// 	//         conn.SendMsg(uint32(msg.MsgID_S2C_BET_NTF), ackData)
// 	//     }
// 	// }
// }
