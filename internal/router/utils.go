package router

import (
	"encoding/json"
	"github.com/aceld/zinx/ziface"
	"xizexcample/internal/logic"
	"xizexcample/internal/msg"
)

// sendErrorResponse 发送错误响应
func sendErrorResponse(conn ziface.IConnection, msgID uint32, errorMsg string) {
	errAck := map[string]interface{}{
		"ret_code": 1,
		"message":  errorMsg,
	}
	ackData, _ := json.Marshal(errAck)
	conn.SendMsg(msgID, ackData)
}

// broadcastRoomState 广播房间状态给所有玩家
func broadcastRoomState(room *logic.Room) {
	ntf := &msg.S2C_SyncRoomStateNtf{
		RoomInfo: buildRoomInfo(room, nil),
	}
	data, _ := json.Marshal(ntf)

	for _, p := range room.GetPlayers() {
		if p.IsOnline() {
			p.Conn.SendMsg(uint32(msg.S2C_SYNC_ROOM_STATE_NTF), data)
		}
	}
}

// sendFullRoomState 向指定玩家发送完整的房间状态（包括手牌）
func sendFullRoomState(player *logic.Player) {
	room := server.GetRoomManager().GetRoomByPlayerID(player.ID)
	if room == nil {
		return
	}

	ntf := &msg.S2C_SyncRoomStateNtf{
		RoomInfo: buildRoomInfo(room, player),
	}
	data, _ := json.Marshal(ntf)
	player.Conn.SendMsg(uint32(msg.S2C_SYNC_ROOM_STATE_NTF), data)
}

// buildRoomInfo 构建房间信息，可选择性地为特定玩家包含手牌信息
func buildRoomInfo(room *logic.Room, targetPlayer *logic.Player) *msg.RoomInfo {
	roomInfo := &msg.RoomInfo{
		RoomId:    uint32(room.ID),
		Players:   []*msg.PlayerInfo{},
		GameState: msg.GameState(room.GetFSM().GetCurrentState()),
		BankerId:  int64(room.GetBankerID()),
	}

	for _, p := range room.GetPlayers() {
		playerInfo := &msg.PlayerInfo{
			PlayerId: uint64(p.ID),
			Nickname: p.Nickname,
			Score:    uint64(p.Score),
			Status:   msg.PlayerStatus(p.GetStatus()),
			IsBanker: p.IsBanker(),
		}
		// 如果是目标玩家，则包含其手牌信息
		if targetPlayer != nil && p.ID == targetPlayer.ID {
			playerInfo.HandCards = cardsToProto(p.Hand)
		}
		roomInfo.Players = append(roomInfo.Players, playerInfo)
	}
	return roomInfo
}

// cardsToProto 将 []logic.Card 转换为 []*msg.Card
func cardsToProto(cards []logic.Card) []*msg.Card {
	protoCards := make([]*msg.Card, len(cards))
	for i, c := range cards {
		protoCards[i] = &msg.Card{
			Suit:   msg.CardSuit(c.Suit),
			Number: msg.CardNumber(c.Number),
		}
	}
	return protoCards
}
