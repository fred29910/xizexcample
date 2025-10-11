package logic

import (
	"testing"
)

func TestRoomAddPlayer(t *testing.T) {
	room := NewRoom(201)
	player1 := NewPlayer(1001, "Player1", nil)
	player2 := NewPlayer(1002, "Player2", nil)

	// 测试添加第一个玩家
	err := room.AddPlayer(player1)
	if err != nil {
		t.Fatalf("AddPlayer failed for player1: %v", err)
	}

	if room.GetPlayerCount() != 1 {
		t.Errorf("Expected 1 player, got %d", room.GetPlayerCount())
	}

	if player1.GetRoomID() != room.ID {
		t.Errorf("Player's room ID is not set correctly. Expected %d, got %d", room.ID, player1.GetRoomID())
	}

	// 测试添加第二个玩家
	err = room.AddPlayer(player2)
	if err != nil {
		t.Fatalf("AddPlayer failed for player2: %v", err)
	}

	if room.GetPlayerCount() != 2 {
		t.Errorf("Expected 2 players, got %d", room.GetPlayerCount())
	}

	// 测试添加重复玩家
	err = room.AddPlayer(player1)
	if err == nil {
		t.Error("Expected error for adding duplicate player, but got nil")
	}

	// 测试房间已满
	player3 := NewPlayer(1003, "Player3", nil)
	player4 := NewPlayer(1004, "Player4", nil)
	player5 := NewPlayer(1005, "Player5", nil)
	player6 := NewPlayer(1006, "Player6", nil)

	room.AddPlayer(player3)
	room.AddPlayer(player4)
	room.AddPlayer(player5)

	err = room.AddPlayer(player6)
	if err == nil {
		t.Error("Expected error for adding player to full room, but got nil")
	}
}

func TestRoomRemovePlayer(t *testing.T) {
	room := NewRoom(202)
	player1 := NewPlayer(1007, "Player1", nil)
	player2 := NewPlayer(1008, "Player2", nil)

	room.AddPlayer(player1)
	room.AddPlayer(player2)

	// 测试移除一个玩家
	err := room.RemovePlayer(player1.ID)
	if err != nil {
		t.Fatalf("RemovePlayer failed: %v", err)
	}

	if room.GetPlayerCount() != 1 {
		t.Errorf("Expected 1 player after removal, got %d", room.GetPlayerCount())
	}

	if player1.GetRoomID() != 0 {
		t.Errorf("Player's room ID should be 0 after removal, got %d", player1.GetRoomID())
	}

	// 测试移除不存在的玩家
	err = room.RemovePlayer(9999)
	if err == nil {
		t.Error("Expected error for removing non-existent player, but got nil")
	}

	// 测试移除最后一个玩家
	err = room.RemovePlayer(player2.ID)
	if err != nil {
		t.Fatalf("RemovePlayer failed for last player: %v", err)
	}

	if room.GetPlayerCount() != 0 {
		t.Errorf("Expected 0 players after removing last player, got %d", room.GetPlayerCount())
	}
}

func TestRoomGetPlayer(t *testing.T) {
	room := NewRoom(203)
	player1 := NewPlayer(1009, "Player1", nil)

	room.AddPlayer(player1)

	// 测试获取存在的玩家
	retrievedPlayer, err := room.GetPlayer(player1.ID)
	if err != nil {
		t.Fatalf("GetPlayer failed: %v", err)
	}

	if retrievedPlayer.ID != player1.ID {
		t.Errorf("Expected player ID %d, got %d", player1.ID, retrievedPlayer.ID)
	}

	// 测试获取不存在的玩家
	_, err = room.GetPlayer(9999)
	if err == nil {
		t.Error("Expected error for getting non-existent player, but got nil")
	}
}
