package router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"../logic"
	"../server"
	"../msg"
	"encoding/json"
)

// BidBankerHandler 处理抢庄请求
type BidBankerHandler struct {
	BaseRouter
}

// Handle 处理请求
func (h *BidBankerHandler) Handle(request ziface.IRequest) {
	// 1. 解析客户端请求
	var bidReq msg.C2S_BidBankerReq
	err := json.Unmarshal(request.GetData(), &bidReq)
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

	// 3. 检查游戏状态是否允许抢庄
	if playerRoom.GetFSM().GetCurrentState() != logic.STATE_BIDDING {
		h.sendErrorResponse(request.GetConnection(), "Cannot bid banker at this time")
		return
	}

	// 4. 处理抢庄逻辑
	// TODO: 实现更复杂的抢庄逻辑，例如比较抢庄倍数
	// 这里简化处理，直接设置第一个抢庄的玩家为庄家
	if !playerRoom.HasBanker() {
		playerRoom.SetBanker(targetPlayer.ID)
		targetPlayer.SetBanker(true)

		// 5. 转换游戏状态到下注阶段
		err = playerRoom.GetFSM().BidBanker()
		if err != nil {
			h.sendErrorResponse(request.GetConnection(), "Failed to transition to betting state")
			return
		}

		// 6. 广播庄家信息
		broadcastBankerInfo(playerRoom)

		// 7. 发送确认响应
		bidAck := &msg.S2C_BidBankerAck{
			RetCode: 0,
			PlayerId: uint64(targetPlayer.ID),
		}
		ackData, err := json.Marshal(bidAck)
		if err != nil {
			h.sendErrorResponse(request.GetConnection(), "Failed to marshal response")
			return
		}
		request.GetConnection().SendMsg(uint32(msg.S2C_BID_BANKER_ACK), ackData)

		// 8. 广播房间状态更新
		broadcastRoomState(playerRoom)
	} else {
		h.sendErrorResponse(request.GetConnection(), "Banker has already been chosen")
	}
}

// broadcastBankerInfo 广播庄家信息给所有玩家
func broadcastBankerInfo(room *logic.Room) {
	bankerInfo := &msg.S2C_BankerNtf{
		BankerId: uint64(room.GetBankerID()),
	}
	ackData, _ := json.Marshal(bankerInfo)

	// TODO: 遍历房间内所有玩家的连接并发送消息
	// for _, player := range room.GetPlayers() {
	//     if conn := getPlayerConnection(player.ID); conn != nil {
	//         conn.SendMsg(uint32(msg.S2C_BANKER_NTF), ackData)
	//     }
	// }
}

// sendErrorResponse 发送错误响应
func (h *BidBankerHandler) sendErrorResponse(conn ziface.IConnection, errorMsg string) {
	errAck := map[string]interface{}{
		"ret_code": 1,
		"message":  errorMsg,
	}
	ackData, _ := json.Marshal(errAck)
	conn.SendMsg(uint32(msg.S2C_BID_BANKER_ACK), ackData)
}