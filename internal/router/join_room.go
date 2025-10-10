package router

import (
	"encoding/json"
	"github.com/aceld/zinx/ziface"
	"xizexcample/internal/logic"
	"xizexcample/internal/msg"
	"xizexcample/internal/server"
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
		sendErrorResponse(request.GetConnection(), uint32(msg.S2C_JOIN_ROOM_ACK), "Invalid request data")
		return
	}

	// 2. 获取或创建房间
	roomManager := server.GetRoomManager()
	room, err := roomManager.GetRoom(int32(joinReq.RoomId))
	if err != nil {
		// 房间不存在，尝试创建
		room, err = roomManager.CreateRoom(int32(joinReq.RoomId))
		if err != nil {
			sendErrorResponse(request.GetConnection(), uint32(msg.S2C_JOIN_ROOM_ACK), "Failed to create room")
			return
		}
	}

	// 3. 检查是否是重连
	playerID := request.GetConnection().GetConnID() // 在实际应用中，应该从请求中获取真实的玩家ID
	existingPlayer, err := room.GetPlayer(int64(playerID))
	if err == nil && !existingPlayer.IsOnline() {
		// 是重连玩家
		// T037: Re-associate connection with the existing Player object
		existingPlayer.Conn = request.GetConnection()
		existingPlayer.SetOnline(true)
		existingPlayer.SetStatus(logic.STATUS_PLAYING) // Or whatever the status was

		// T038: Send a full room state sync message to the reconnected player
		sendFullRoomState(existingPlayer)

		// Broadcast to other players that this player has reconnected
		broadcastRoomState(room)
		return
	}

	// 4. 创建新玩家对象
	player := logic.NewPlayer(int64(playerID), "Player"+string(playerID), request.GetConnection())

	// 5. 将玩家加入房间
	err = room.AddPlayer(player)
	if err != nil {
		sendErrorResponse(request.GetConnection(), uint32(msg.S2C_JOIN_ROOM_ACK), err.Error())
		return
	}

	// 将 playerID 设置到连接属性中
	request.GetConnection().SetProperty("playerID", player.ID)
	roomManager.RegisterPlayer(player.ID, room.ID)

	// 6. 准备并发送成功响应
	joinAck := &msg.S2C_JoinRoomAck{
		RetCode: 0,
		RoomInfo: &msg.RoomInfo{
			RoomId:    uint32(room.ID),
			Players:   []*msg.PlayerInfo{},
			GameState: msg.GameState(room.GetFSM().GetCurrentState()),
			BankerId:  int64(room.GetBankerID()),
		},
	}

	for _, p := range room.GetPlayers() {
		playerInfo := &msg.PlayerInfo{
			PlayerId: uint64(p.ID),
			Nickname: p.Nickname,
			Score:    uint64(p.Score),
			Status:   msg.PlayerStatus(p.GetStatus()),
			IsBanker: p.IsBanker(),
		}
		joinAck.RoomInfo.Players = append(joinAck.RoomInfo.Players, playerInfo)
	}

	ackData, err := json.Marshal(joinAck)
	if err != nil {
		sendErrorResponse(request.GetConnection(), uint32(msg.S2C_JOIN_ROOM_ACK), "Failed to marshal response")
		return
	}

	request.GetConnection().SendMsg(uint32(msg.S2C_JOIN_ROOM_ACK), ackData)

	// 7. 广播房间状态更新给所有玩家
	broadcastRoomState(room)
}

