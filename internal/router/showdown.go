package router

import (
	"encoding/json"
	"github.com/aceld/zinx/ziface"
	"xizexcample/internal/logic"
	"xizexcample/internal/msg"
	"xizexcample/internal/pkg/logger"
)

// ShowdownHandler 处理摊牌请求
type ShowdownHandler struct {
	BaseRouter
}

// Handle 处理请求
func (h *ShowdownHandler) Handle(request ziface.IRequest) {
	// 1. 解析客户端请求
	var showdownReq msg.C2S_ShowdownReq
	err := json.Unmarshal(request.GetData(), &showdownReq)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal showdown request: %v", err)
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_SHOWDOWN_ACK), "Invalid request data")
		return
	}

	// 2. 获取玩家和房间
	player, room, err := GetPlayerAndRoom(request)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get player and room: %v", err)
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_SHOWDOWN_ACK), err.Error())
		return
	}

	// 3. 检查游戏状态
	if room.GetFSM().GetCurrentState() != logic.STATE_SHOWDOWN {
		sendErrorResponse(request.GetConnection(), uint32(msg.MsgID_S2C_SHOWDOWN_ACK), "Cannot showdown at this time")
		return
	}

	// 4. 计算牌型和分数 (简化逻辑)
	// 在实际应用中，这里应该有复杂的牌型计算逻辑
	// player.CalculateHand()
	logger.InfoLogger.Printf("Player %d in room %d shows hand", player.ID, room.ID)

	// 5. 检查是否所有人都已摊牌
	allShowdown := true
	for _, p := range room.GetPlayers() {
		// 简化检查逻辑，实际应检查玩家是否已摊牌的状态
		if p.GetStatus() == logic.STATUS_PLAYING {
			// allShowdown = false
			// break
		}
	}

	// 6. 如果所有人都已摊牌，进入结算状态
	if allShowdown {
		err := room.GetFSM().Showdown()
		if err != nil {
			logger.ErrorLogger.Printf("Failed to transition to settlement state in room %d: %v", room.ID, err)
		}
	}

	// 7. 发送确认响应
	ack := &msg.S2C_ShowdownAck{RetCode: 0}
	ackData, _ := json.Marshal(ack)
	request.GetConnection().SendMsg(uint32(msg.MsgID_S2C_SHOWDOWN_ACK), ackData)

	// 8. 广播房间状态
	broadcastRoomState(room)
}
