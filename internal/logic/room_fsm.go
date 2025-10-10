package logic

import (
	"errors"
	"fmt"
)

// GameState 游戏状态
type GameState int

const (
	STATE_UNKNOWN GameState = iota
	STATE_WAITING_FOR_PLAYERS
	STATE_DEALING
	STATE_BIDDING
	STATE_BETTING
	STATE_SHOWDOWN
	STATE_SETTLEMENT
)

// RoomFSM 房间状态机
type RoomFSM struct {
	currentState GameState
	room         *Room
}

// NewRoomFSM 创建一个新的房间状态机
func NewRoomFSM(room *Room) *RoomFSM {
	return &RoomFSM{
		currentState: STATE_WAITING_FOR_PLAYERS,
		room:         room,
	}
}

// GetCurrentState 获取当前状态
func (fsm *RoomFSM) GetCurrentState() GameState {
	return fsm.currentState
}

// TransitionTo 转换到新状态
func (fsm *RoomFSM) TransitionTo(newState GameState) error {
	// 定义状态转换规则
	validTransitions := map[GameState][]GameState{
		STATE_WAITING_FOR_PLAYERS: {STATE_DEALING},
		STATE_DEALING:            {STATE_BIDDING},
		STATE_BIDDING:            {STATE_BETTING},
		STATE_BETTING:            {STATE_SHOWDOWN},
		STATE_SHOWDOWN:           {STATE_SETTLEMENT},
		STATE_SETTLEMENT:         {STATE_WAITING_FOR_PLAYERS}, // 一局结束后，等待下一局
	}

	// 检查转换是否有效
	if allowedStates, exists := validTransitions[fsm.currentState]; exists {
		for _, allowedState := range allowedStates {
			if newState == allowedState {
				fsm.currentState = newState
				return nil
			}
		}
	}

	return fmt.Errorf("invalid state transition from %d to %d", fsm.currentState, newState)
}

// CanStartGame 检查是否可以开始游戏
func (fsm *RoomFSM) CanStartGame() bool {
	return fsm.currentState == STATE_WAITING_FOR_PLAYERS && fsm.room.GetPlayerCount() >= 2
}

// StartGame 开始游戏，转换到发牌状态
func (fsm *RoomFSM) StartGame() error {
	if !fsm.CanStartGame() {
		return errors.New("cannot start game at this time")
	}

	// TODO: 广播 S2C_GameStartNtf 给所有玩家
	// broadcastGameStart(fsm.room)

	return fsm.TransitionTo(STATE_DEALING)
}

// DealCards 发牌
func (fsm *RoomFSM) DealCards() error {
	if fsm.currentState != STATE_DEALING {
		return errors.New("cannot deal cards in current state")
	}

	// 重置并洗牌
	fsm.room.ResetDeck()

	// 给每个玩家发5张牌
	for _, player := range fsm.room.GetPlayers() {
		err := fsm.room.DealCardsToPlayer(player.ID, 5)
		if err != nil {
			// 记录错误，但继续给其他玩家发牌
			// TODO: 添加日志
			continue
		}
		// 更新玩家状态为 PLAYING
		player.SetStatus(STATUS_PLAYING)

		// TODO: 广播 S2C_DealCardsNtf 给该玩家
		// broadcastDealCards(player, player.GetHand())
	}

	// 发牌完成后，转换到抢庄状态
	return fsm.TransitionTo(STATE_BIDDING)
}

// BidBanker 抢庄
func (fsm *RoomFSM) BidBanker() error {
	if fsm.currentState != STATE_BIDDING {
		return errors.New("cannot bid banker in current state")
	}
	// TODO: 实现抢庄逻辑，例如选择抢庄倍数最高的玩家
	// 这里简化处理，直接转换到下注状态
	return fsm.TransitionTo(STATE_BETTING)
}

// PlaceBet 下注
func (fsm *RoomFSM) PlaceBet() error {
	if fsm.currentState != STATE_BETTING {
		return errors.New("cannot place bet in current state")
	}
	// TODO: 实现下注逻辑
	// 这里简化处理，直接转换到摊牌状态
	return fsm.TransitionTo(STATE_SHOWDOWN)
}

// Showdown 摊牌
func (fsm *RoomFSM) Showdown() error {
	if fsm.currentState != STATE_SHOWDOWN {
		return errors.New("cannot showdown in current state")
	}
	// TODO: 实现摊牌和结算逻辑
	// 这里简化处理，直接转换到结算状态
	return fsm.TransitionTo(STATE_SETTLEMENT)
}

// Settlement 结算
func (fsm *RoomFSM) Settlement() error {
	if fsm.currentState != STATE_SETTLEMENT {
		return errors.New("cannot settlement in current state")
	}
	// TODO: 实现结算逻辑，更新玩家分数
	// 结算完成后，转换到等待玩家状态，准备下一局
	return fsm.TransitionTo(STATE_WAITING_FOR_PLAYERS)
}
