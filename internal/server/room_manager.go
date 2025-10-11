package server

import (
	"errors"
	"sync"
	"xizexcample/internal/logic"
)

// RoomManager 是一个全局单例，用于管理所有房间
var (
	roomManagerInstance *RoomManager
	once                sync.Once
)

// RoomManager 房间管理器
type RoomManager struct {
	rooms      map[int32]*logic.Room // key: roomID
	playerRoom map[int64]int32       // key: playerID, value: roomID
	mu         sync.RWMutex
}

// GetRoomManager 获取 RoomManager 单例
func GetRoomManager() *RoomManager {
	once.Do(func() {
		roomManagerInstance = &RoomManager{
			rooms:      make(map[int32]*logic.Room),
			playerRoom: make(map[int64]int32),
		}
	})
	return roomManagerInstance
}

// CreateRoom 创建一个新房间
func (rm *RoomManager) CreateRoom(roomID int32) (*logic.Room, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomID]; exists {
		return nil, errors.New("room already exists")
	}

	room := logic.NewRoom(roomID)
	rm.rooms[roomID] = room
	return room, nil
}

// GetRoom 根据房间ID获取房间
func (rm *RoomManager) GetRoom(roomID int32) (*logic.Room, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

// DeleteRoom 删除房间
func (rm *RoomManager) DeleteRoom(roomID int32) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomID]; !exists {
		return errors.New("room not found")
	}

	delete(rm.rooms, roomID)
	return nil
}

// GetRoomByPlayerID 根据玩家ID获取房间
func (rm *RoomManager) GetRoomByPlayerID(playerID int64) *logic.Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	roomID, ok := rm.playerRoom[playerID]
	if !ok {
		return nil
	}

	room, ok := rm.rooms[roomID]
	if !ok {
		// Data inconsistency, should not happen
		return nil
	}
	return room
}

// RegisterPlayer 将玩家注册到房间
func (rm *RoomManager) RegisterPlayer(playerID int64, roomID int32) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.playerRoom[playerID] = roomID
}

// UnregisterPlayer 从房间注销玩家
func (rm *RoomManager) UnregisterPlayer(playerID int64) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.playerRoom, playerID)
}
