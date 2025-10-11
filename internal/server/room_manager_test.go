package server

import (
	"testing"
)

func TestRoomManagerCreateAndGetRoom(t *testing.T) {
	rm := GetRoomManager()

	// 测试创建房间
	roomID := int32(101)
	room, err := rm.CreateRoom(roomID)
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}

	if room.ID != roomID {
		t.Errorf("Expected room ID %d, got %d", roomID, room.ID)
	}

	// 测试获取已存在的房间
	retrievedRoom, err := rm.GetRoom(roomID)
	if err != nil {
		t.Fatalf("GetRoom failed: %v", err)
	}

	if retrievedRoom.ID != roomID {
		t.Errorf("Expected room ID %d, got %d", roomID, retrievedRoom.ID)
	}

	// 测试获取不存在的房间
	_, err = rm.GetRoom(999)
	if err == nil {
		t.Error("Expected error for non-existent room, but got nil")
	}
}

func TestRoomManagerCreateDuplicateRoom(t *testing.T) {
	rm := GetRoomManager()

	roomID := int32(102)
	_, err := rm.CreateRoom(roomID)
	if err != nil {
		t.Fatalf("First CreateRoom failed: %v", err)
	}

	// 尝试创建相同ID的房间
	_, err = rm.CreateRoom(roomID)
	if err == nil {
		t.Error("Expected error for creating duplicate room, but got nil")
	}
}

func TestRoomManagerDeleteRoom(t *testing.T) {
	rm := GetRoomManager()

	roomID := int32(103)
	_, err := rm.CreateRoom(roomID)
	if err != nil {
		t.Fatalf("CreateRoom failed: %v", err)
	}

	err = rm.DeleteRoom(roomID)
	if err != nil {
		t.Fatalf("DeleteRoom failed: %v", err)
	}

	// 验证房间已被删除
	_, err = rm.GetRoom(roomID)
	if err == nil {
		t.Error("Expected error for getting deleted room, but got nil")
	}
}
