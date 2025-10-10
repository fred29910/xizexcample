package router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"xizexcample/internal/logic"
	"xizexcample/internal/msg"
	"xizexcample/internal/server"
	"encoding/json"
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
		h.sendErrorResponse(request.GetConnection(), "Invalid request data")
		return
	}

	// 2. 获取玩家和房间 (这里简化处理)
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

	// 3. 检查游戏状态是否允许下注
	if playerRoom.GetFSM().GetCurrentState() != logic.STATE_BETTING {
		h.sendErrorResponse(request.GetConnection(), "Cannot place bet at this time")
		return
	}

	// 4. 检查玩家是否已经下注
	if targetPlayer.HasBet() {
		h.sendErrorResponse(request.GetConnection(), "Player has already placed a bet")
		return
	}

	// 5. 检查下注金额是否有效
	if betReq.BetAmount <= 0 {
		h.sendErrorResponse(request.GetConnection(), "Invalid bet amount")
		return
	}

	// TODO: 检查玩家余额是否足够
	// if targetPlayer.GetScore() < betReq.BetAmount {
	//     h.sendErrorResponse(request.GetConnection(), "Insufficient balance")
	//     return
	// }

	// 6. 处理下注逻辑
	targetPlayer.PlaceBet(betReq.BetAmount)

	// 7. 广播下注信息
	broadcastBetInfo(playerRoom, targetPlayer, betReq.BetAmount)

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
		err = playerRoom.GetFSM().PlaceBet()
		if err != nil {
			h.sendErrorResponse(request.GetConnection(), "Failed to transition to showdown state")
			return
		}

		// TODO: 自动触发摊牌逻辑
		// go h.triggerShowdown(playerRoom)
	}

	// 10. 发送确认响应
	betAck := &msg.S2C_PlaceBetAck{
		RetCode: 0,
		BetAmount: betReq.BetAmount,
	}
	ackData, err := json.Marshal(betAck)
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), "Failed to marshal response")
		return
	}
	request.GetConnection().SendMsg(uint32(msg.S2C_PLACE_BET_ACK), ackData)

	// 11. 广播房间状态更新
	broadcastRoomState(playerRoom)
}

// broadcastBetInfo 广播下注信息给所有玩家
func broadcastBetInfo(room *logic.Room, player *logic.Player, amount uint64) {
	betInfo := &msg.S2C_BetNtf{
		PlayerId: uint64(player.ID),
		BetAmount: amount,
	}
	ackData, _ := json.Marshal(betInfo)

	// TODO: 遍历房间内所有玩家的连接并发送消息
	// for _, p := range room.GetPlayers() {
	//     if conn := getPlayerConnection(p.ID); conn != nil {
	//         conn.SendMsg(uint32(msg.S2C_BET_NTF), ackData)
	//     }
	// }
}

// sendErrorResponse 发送错误响应
func (h *PlaceBetHandler) sendErrorResponse(conn ziface.IConnection, errorMsg string) {
	errAck := map[string]interface{}{
		"ret_code": 1,
		"message":  errorMsg,
	}
	ackData, _ := json.Marshal(errAck)
	conn.SendMsg(uint32(msg.S2C_PLACE_BET_ACK), ackData)
}
