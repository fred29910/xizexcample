package router

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"../logic"
	"../server"
	"../msg"
	"encoding/json"
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
		h.sendErrorResponse(request.GetConnection(), "Invalid request data")
		return
	}

	// 2. 获取玩家 (这里简化处理，实际应该从连接中获取玩家)
	playerID := request.GetConnection().GetConnID()
	roomManager := server.GetRoomManager()

	// 假设玩家所在的房间ID是已知的，或者从连接上下文中获取
	// 这里我们遍历所有房间来找到玩家，这是一个简化的实现
	var playerRoom *logic.Room
	var targetPlayer *logic.Player
	for _, room := range roomManager.GetAllRooms() { // 假设 RoomManager 有这个方法
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

	// 3. 更新玩家状态
	if readyReq.IsReady {
		targetPlayer.SetStatus(logic.STATUS_READY)
	} else {
		targetPlayer.SetStatus(logic.STATUS_WAITING)
	}

	// 4. 检查是否可以开始游戏
	if canStartGame(playerRoom) {
		startGame(playerRoom)
	}

	// 5. 广播房间状态更新给所有玩家
	broadcastRoomState(playerRoom)

	// 6. 发送确认响应给当前玩家
	readyAck := &msg.S2C_SyncRoomStateNtf{
		RoomInfo: createRoomInfo(playerRoom),
	}
	ackData, err := json.Marshal(readyAck)
	if err != nil {
		h.sendErrorResponse(request.GetConnection(), "Failed to marshal response")
		return
	}
	request.GetConnection().SendMsg(uint32(msg.S2C_SYNC_ROOM_STATE_NTF), ackData)
}

// canStartGame 检查游戏是否可以开始
func canStartGame(room *logic.Room) bool {
	if room == nil {
		return false
	}
	playerCount := room.GetPlayerCount()
	return playerCount >= 2 && playerCount <= 5
}

// startGame 开始游戏
func startGame(room *logic.Room) {
	// TODO: 实现游戏开始的逻辑
	// 1. 重置牌堆
	room.ResetDeck()
	// 2. 给每个玩家发牌
	for _, player := range room.GetPlayers() {
		err := room.DealCardsToPlayer(player.ID, 5)
		if err != nil {
			// 处理错误，例如日志记录
			continue
		}
		player.SetStatus(logic.STATUS_PLAYING)
	}
	// 3. 更新游戏状态
	// TODO: room.SetGameState(logic.GAME_STATE_DEALING)
}

// createRoomInfo 创建 RoomInfo 消息
func createRoomInfo(room *logic.Room) *msg.RoomInfo {
	roomInfo := &msg.RoomInfo{
		RoomId:     uint32(room.ID),
		GameState: msg.GAME_STATE_WAITING_FOR_PLAYERS, // TODO: 从 room 获取真实状态
		BankerId:   0,
		Players:    []*msg.PlayerInfo{},
	}

	for _, p := range room.GetPlayers() {
		playerInfo := &msg.PlayerInfo{
			PlayerId: uint64(p.ID),
			Nickname: p.Nickname,
			Score:    uint64(p.Score),
			Status:   msg.PlayerStatus(p.GetStatus()),
			IsBanker: p.IsBanker,
		}
		roomInfo.Players = append(roomInfo.Players, playerInfo)
	}

	return roomInfo
}

// broadcastRoomState 广播房间状态给所有玩家
func broadcastRoomState(room *logic.Room) {
	roomInfo := createRoomInfo(room)
	ackData, _ := json.Marshal(roomInfo)

	// TODO: 遍历房间内所有玩家的连接并发送消息
	// for _, player := range room.GetPlayers() {
	//     if conn := getPlayerConnection(player.ID); conn != nil {
	//         conn.SendMsg(uint32(msg.S2C_SYNC_ROOM_STATE_NTF), ackData)
	//     }
	// }
}

// sendErrorResponse 发送错误响应
func (h *PlayerReadyHandler) sendErrorResponse(conn ziface.IConnection, errorMsg string) {
	errAck := map[string]interface{}{
		"ret_code": 1,
		"message":  errorMsg,
	}
	ackData, _ := json.Marshal(errAck)
	conn.SendMsg(uint32(msg.S2C_SYNC_ROOM_STATE_NTF), ackData)
}