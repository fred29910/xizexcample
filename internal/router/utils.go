package router

import (
	"encoding/json"
	"errors"
	"github.com/aceld/zinx/ziface"
	"xizexcample/internal/logic"
	"xizexcample/internal/msg"
	"xizexcample/internal/server"
)

// sendErrorResponse 向客户端发送一个标准格式的错误响应
func sendErrorResponse(conn ziface.IConnection, msgID uint32, errorMsg string) {
	errAck := map[string]interface{}{
		"ret_code": 1,
		"message":  errorMsg,
	}
	ackData, _ := json.Marshal(errAck)
	conn.SendMsg(msgID, ackData)
}

// GetPlayerAndRoom 从连接中获取玩家和房间
func GetPlayerAndRoom(request ziface.IRequest) (*logic.Player, *logic.Room, error) {
	playerID, err := request.GetConnection().GetProperty("playerID")
	if err != nil {
		return nil, nil, errors.New("player not logged in")
	}

	room := server.GetRoomManager().GetRoomByPlayerID(playerID.(int64))
	if room == nil {
		return nil, nil, errors.New("player not in any room")
	}

	player, err := room.GetPlayer(playerID.(int64))
	if err != nil {
		return nil, nil, errors.New("player not found in room")
	}

	return player, room, nil
}

// broadcastRoomState 广播房间状态给所有玩家
func broadcastRoomState(room *logic.Room) {
	players := room.GetPlayers()
	playerInfos := make([]*msg.PlayerInfo, len(players))
	for i, p := range players {
		playerInfos[i] = &msg.PlayerInfo{
			PlayerId: p.ID,
			Nickname: p.Nickname,
			Score:    p.Score,
			Status:   msg.PlayerStatus(p.GetStatus()),
			IsBanker: p.IsBanker(),
		}
	}

	roomInfo := &msg.RoomInfo{
		RoomId:    room.ID,
		Players:   playerInfos,
		GameState: msg.GameState(room.GetFSM().GetCurrentState()),
		BankerId:  room.GetBankerID(),
	}

	ntf := &msg.S2C_SyncRoomStateNtf{
		RoomInfo: roomInfo,
	}
	ntfData, _ := json.Marshal(ntf)

	for _, p := range players {
		if p.Conn != nil && p.IsOnline() {
			p.Conn.SendMsg(uint32(msg.MsgID_S2C_SYNC_ROOM_STATE_NTF), ntfData)
		}
	}
}

// sendFullRoomState 向单个玩家发送完整的房间状态
func sendFullRoomState(player *logic.Player) {
	room := server.GetRoomManager().GetRoomByPlayerID(player.ID)
	if room == nil {
		return
	}

	players := room.GetPlayers()
	playerInfos := make([]*msg.PlayerInfo, len(players))
	for i, p := range players {
		playerInfos[i] = &msg.PlayerInfo{
			PlayerId: p.ID,
			Nickname: p.Nickname,
			Score:    p.Score,
			Status:   msg.PlayerStatus(p.GetStatus()),
			IsBanker: p.IsBanker(),
		}
	}

	roomInfo := &msg.RoomInfo{
		RoomId:    room.ID,
		Players:   playerInfos,
		GameState: msg.GameState(room.GetFSM().GetCurrentState()),
		BankerId:  room.GetBankerID(),
	}

	ntf := &msg.S2C_SyncRoomStateNtf{
		RoomInfo: roomInfo,
	}
	ntfData, _ := json.Marshal(ntf)

	if player.Conn != nil && player.IsOnline() {
		player.Conn.SendMsg(uint32(msg.MsgID_S2C_SYNC_ROOM_STATE_NTF), ntfData)
	}
}
