package server

import (
	"errors"
	"sync"
)

// RoomManager 是一个全局单例，用于管理所有房间
var (
	roomManagerInstance *RoomManager
	once                sync.Once
)

// RoomManager 房间管理器
type RoomManager struct {
	rooms map[int32]*Room // key: roomID
	mu    sync.RWMutex
}

// GetRoomManager 获取 RoomManager 单例
func GetRoomManager() *RoomManager {
	once.Do(func() {
		roomManagerInstance = &RoomManager{
			rooms: make(map[int32]*Room),
		}
	})
	return roomManagerInstance
}

// CreateRoom 创建一个新房间
func (rm *RoomManager) CreateRoom(roomID int32) (*Room, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, exists := rm.rooms[roomID]; exists {
		return nil, errors.New("room already exists")
	}

	room := NewRoom(roomID)
	rm.rooms[roomID] = room
	return room, nil
}

// GetRoom 根据房间ID获取房间
func (rm *RoomManager) GetRoom(roomID int32) (*Room, error) {
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

// Room 房间结构体 (简化版，将在后续任务中完善)
type Room struct {
	ID      int32
	Players map[int64]*Player // key: playerID
}

// NewRoom 创建一个新房间实例
func NewRoom(roomID int32) *Room {
	return &Room{
		ID:      roomID,
		Players: make(map[int64]*Player),
	}
}

// Player 玩家结构体 (简化版，将在后续任务中完善)
type Player struct {
	ID   int64
	Name string
}