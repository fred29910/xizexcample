package logic

import (
	"testing"
)

func TestRoomFSMBidBanker(t *testing.T) {
	room := NewRoom(301)
	fsm := NewRoomFSM(room)

	// 测试从无效状态抢庄
	err := fsm.BidBanker()
	if err == nil {
		t.Error("Expected error when bidding banker from invalid state, but got nil")
	}

	// 测试从抢庄状态抢庄
	err = fsm.TransitionTo(STATE_BIDDING)
	if err != nil {
		t.Fatalf("Failed to transition to BIDDING state: %v", err)
	}

	err = fsm.BidBanker()
	if err != nil {
		t.Errorf("BidBanker failed from BIDDING state: %v", err)
	}

	// 验证状态已转换
	if fsm.GetCurrentState() != STATE_BETTING {
		t.Errorf("Expected state to be BETTING after BidBanker, got %d", fsm.GetCurrentState())
	}
}

func TestRoomFSMPlaceBet(t *testing.T) {
	room := NewRoom(302)
	fsm := NewRoomFSM(room)

	// 测试从无效状态下注
	err := fsm.PlaceBet()
	if err == nil {
		t.Error("Expected error when placing bet from invalid state, but got nil")
	}

	// 测试从下注状态下注
	err = fsm.TransitionTo(STATE_BETTING)
	if err != nil {
		t.Fatalf("Failed to transition to BETTING state: %v", err)
	}

	err = fsm.PlaceBet()
	if err != nil {
		t.Errorf("PlaceBet failed from BETTING state: %v", err)
	}

	// 验证状态已转换
	if fsm.GetCurrentState() != STATE_SHOWDOWN {
		t.Errorf("Expected state to be SHOWDOWN after PlaceBet, got %d", fsm.GetCurrentState())
	}
}

func TestRoomFSMShowdown(t *testing.T) {
	room := NewRoom(303)
	fsm := NewRoomFSM(room)

	// 测试从无效状态摊牌
	err := fsm.Showdown()
	if err == nil {
		t.Error("Expected error when showdown from invalid state, but got nil")
	}

	// 测试从摊牌状态摊牌
	err = fsm.TransitionTo(STATE_SHOWDOWN)
	if err != nil {
		t.Fatalf("Failed to transition to SHOWDOWN state: %v", err)
	}

	err = fsm.Showdown()
	if err != nil {
		t.Errorf("Showdown failed from SHOWDOWN state: %v", err)
	}

	// 验证状态已转换
	if fsm.GetCurrentState() != STATE_SETTLEMENT {
		t.Errorf("Expected state to be SETTLEMENT after Showdown, got %d", fsm.GetCurrentState())
	}
}

func TestRoomFSMSettlement(t *testing.T) {
	room := NewRoom(304)
	fsm := NewRoomFSM(room)

	// 测试从无效状态结算
	err := fsm.Settlement()
	if err == nil {
		t.Error("Expected error when settlement from invalid state, but got nil")
	}

	// 测试从结算状态结算
	err = fsm.TransitionTo(STATE_SETTLEMENT)
	if err != nil {
		t.Fatalf("Failed to transition to SETTLEMENT state: %v", err)
	}

	err = fsm.Settlement()
	if err != nil {
		t.Errorf("Settlement failed from SETTLEMENT state: %v", err)
	}

	// 验证状态已转换
	if fsm.GetCurrentState() != STATE_WAITING_FOR_PLAYERS {
		t.Errorf("Expected state to be WAITING_FOR_PLAYERS after Settlement, got %d", fsm.GetCurrentState())
	}
}

func TestRoomFSMInvalidTransition(t *testing.T) {
	room := NewRoom(305)
	fsm := NewRoomFSM(room)

	// 测试从 WAITING_FOR_PLAYERS 直接转换到 BETTING
	err := fsm.TransitionTo(STATE_BETTING)
	if err == nil {
		t.Error("Expected error for invalid transition, but got nil")
	}
}