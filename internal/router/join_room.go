package router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"xizexcample/internal/logic"
	"xizexcample/internal/server"
	"xizexcample/internal/msg" // 假设这是生成的 protobuf 消息包
	"encoding/json"
)

// JoinRoomHandler 处理加入房间请求
type JoinRoomHandler struct {
	BaseRouter
}

// Handle 处理请求
func (h *JoinRoomHandler) Handle(request ziface.IRequest) {
	// 1. 解析客户端请求
	var joinReq msg.C2S_JoinRoomReq
	err := json.Unmarshal(request.GetData(), &joinReq)
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), "Invalid request data")
		return
	}

	// 2. 获取或创建房间
	roomManager := server.GetRoomManager()
	room, err := roomManager.GetRoom(joinReq.RoomId)
	if err != nil {
		// 房间不存在，尝试创建
		room, err = roomManager.CreateRoom(joinReq.RoomId)
		if err != nil {
			h.sendErrorResponse(request.GetConnection(), "Failed to create room")
			return
		}
	}

	// 3. 创建玩家对象 (这里简化处理，实际应该从连接中获取或创建玩家)
	// TODO: 需要将 Zinx 连接与 Player 对象关联起来
	playerID := request.GetConnection().GetConnID() // 使用连接ID作为临时玩家ID
	player := logic.NewPlayer(playerID, "Player"+string(playerID))

	// 4. 将玩家加入房间
	err = room.AddPlayer(player)
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), err.Error())
		return
	}

	// 5. 准备并发送成功响应
	joinAck := &msg.S2C_JoinRoomAck{
		RetCode: 0,
		RoomInfo: &msg.RoomInfo{
			RoomId: uint32(room.ID),
			Players: []*msg.PlayerInfo{},
			GameState: msg.GAME_STATE_WAITING_FOR_PLAYERS,
			BankerId: 0,
		},
	}

	// 将 RoomInfo 中的 Players 填充
	for _, p := range room.GetPlayers() {
		playerInfo := &msg.PlayerInfo{
			PlayerId: uint64(p.ID),
			Nickname: p.Nickname,
			Score:    uint64(p.Score),
			Status:   msg.PLAYER_STATUS_WAITING,
			IsBanker: false,
		}
		joinAck.RoomInfo.Players = append(joinAck.RoomInfo.Players, playerInfo)
	}

	ackData, err := json.Marshal(joinAck)
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), "Failed to marshal response")
		return
	}

	request.GetConnection().SendMsg(uint32(msg.S2C_JOIN_ROOM_ACK), ackData)

	// 6. 广播房间状态更新给所有玩家
	broadcastRoomState(room)
}

// sendErrorResponse 发送错误响应
func (h *JoinRoomHandler) sendErrorResponse(conn ziface.IConnection, errorMsg string) {
	// TODO: 定义一个标准的错误消息结构
	errAck := map[string]interface{}{
		"ret_code": 1,
		"message":  errorMsg,
	}
	ackData, _ := json.Marshal(errAck)
	conn.SendMsg(uint32(msg.S2C_JOIN_ROOM_ACK), ackData)
}

// broadcastRoomState 广播房间状态给所有玩家
func broadcastRoomState(room *logic.Room) {
	// TODO: 实现广播逻辑，需要遍历房间内所有玩家的连接并发送消息
	// 这部分逻辑需要将 Player 对象与 Zinx Connection 对象关联起来
	// 目前只是一个占位符
	roomInfo := &msg.RoomInfo{
		RoomId: uint32(room.ID),
		Players: []*msg.PlayerInfo{},
		GameState: msg.GAME_STATE_WAITING_FOR_PLAYERS,
		BankerId: 0,
	}

	for _, p := range room.GetPlayers() {
		playerInfo := &msg.PlayerInfo{
			PlayerId: uint64(p.ID),
			Nickname: p.Nickname,
			Score:    uint64(p.Score),
			Status:   msg.PLAYER_STATUS_WAITING,
			IsBanker: false,
		}
		roomInfo.Players = append(roomInfo.Players, playerInfo)
	}

	ackData, _ := json.Marshal(roomInfo)
	// TODO: conn.SendMsg(uint32(msg.S2C_SYNC_ROOM_STATE_NTF), ackData)
}
