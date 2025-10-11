package router

import (
	"encoding/json"
	"github.com/aceld/zinx/ziface"
	"xizexcample/internal/logic"
	"xizexcample/internal/msg"
	"xizexcample/internal/pkg/logger"
)

// PlayerReadyHandler 处理玩家准备请求
type PlayerReadyHandler struct {
	BaseRouter
}

// Handle 处理请求
func (h *PlayerReadyHandler) Handle(request ziface.IRequest) {
	// 1. 解析客户端请求
	var readyReq msg.C2S_PlayerReadyReq
	err := json.Unmarshal(request.GetData(), &readyReq)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal player ready request: %v", err)
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_SYNC_ROOM_STATE_NTF), "Invalid request data")
		return
	}

	// 2. 获取玩家和房间
	player, room, err := GetPlayerAndRoom(request)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get player and room: %v", err)
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_SYNC_ROOM_STATE_NTF), err.Error())
		return
	}

	// 3. 检查游戏状态
	if room.GetFSM().GetCurrentState() != logic.STATE_WAITING_FOR_PLAYERS {
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_SYNC_ROOM_STATE_NTF), "Game has already started")
		return
	}

	// 4. 设置玩家状态
	if readyReq.IsReady {
		player.SetStatus(logic.STATUS_READY)
	} else {
		player.SetStatus(logic.STATUS_WAITING)
	}
	logger.InfoLogger.Printf("Player %d in room %d set status to %v", player.ID, room.ID, player.GetStatus())

	// 5. 检查是否所有玩家都已准备好开始游戏
	allReady := true
	for _, p := range room.GetPlayers() {
		if p.GetStatus() != logic.STATUS_READY {
			allReady = false
			break
		}
	}

	// 6. 如果满足开始条件，则开始游戏
	if allReady && room.GetFSM().CanStartGame() {
		logger.InfoLogger.Printf("All players are ready. Starting game in room %d", room.ID)
		err := room.GetFSM().StartGame()
		if err != nil {
			logger.ErrorLogger.Printf("Failed to start game in room %d: %v", room.ID, err)
		}
	}

	// 7. 广播房间状态更新
	broadcastRoomState(room)
}
