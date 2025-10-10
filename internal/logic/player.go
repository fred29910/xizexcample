package logic

import (
	"sync"
)

// Player 表示一个玩家
type Player struct {
	ID       int64
	Nickname string
	Score    int64
	RoomID   int32 // 0 表示不在任何房间中

	// 游戏相关状态
	Hand      []Card
	BetAmount int32
	Status    PlayerStatus
	IsBanker  bool

	// 连接相关
	// TODO: 将在后续任务中添加 Zinx 连接对象
	// conn znet.Connection

	mu sync.RWMutex
}

// NewPlayer 创建一个新玩家
func NewPlayer(id int64, nickname string) *Player {
	return &Player{
		ID:       id,
		Nickname: nickname,
		Score:    1000, // 初始分数
		RoomID:   0,
		Hand:     make([]Card, 0),
		Status:   STATUS_WAITING,
	}
}

// SetStatus 设置玩家状态
func (p *Player) SetStatus(status PlayerStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Status = status
}

// GetStatus 获取玩家状态
func (p *Player) GetStatus() PlayerStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Status
}

// AddCard 给玩家发一张牌
func (p *Player) AddCard(card Card) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Hand = append(p.Hand, card)
}

// ClearHand 清空手牌
func (p *Player) ClearHand() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Hand = make([]Card, 0)
}

// SetRoomID 设置玩家所在的房间ID
func (p *Player) SetRoomID(roomID int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.RoomID = roomID
}

// GetRoomID 获取玩家所在的房间ID
func (p *Player) GetRoomID() int32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.RoomID
}

// Card 表示一张扑克牌
type Card struct {
	Suit Suit
	Rank Rank
}

// Suit 花色
type Suit int

const (
	SUIT_UNKNOWN Suit = iota
	SUIT_SPADES
	SUIT_HEARTS
	SUIT_CLUBS
	SUIT_DIAMONDS
)

// Rank 点数
type Rank int

const (
	RANK_UNKNOWN Rank = iota
	RANK_ACE
	RANK_TWO
	RANK_THREE
	RANK_FOUR
	RANK_FIVE
	RANK_SIX
	RANK_SEVEN
	RANK_EIGHT
	RANK_NINE
	RANK_TEN
	RANK_JACK
	RANK_QUEEN
	RANK_KING
)

// PlayerStatus 玩家状态
type PlayerStatus int

const (
	STATUS_UNKNOWN PlayerStatus = iota
	STATUS_WAITING
	STATUS_READY
	STATUS_PLAYING
)