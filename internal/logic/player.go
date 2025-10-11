package logic

import (
	"github.com/aceld/zinx/ziface"
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
	isBanker  bool

	// 连接相关
	isOnline       bool
	Conn           ziface.IConnection
	DisconnectTime int64 // Unix timestamp

	mu sync.RWMutex
}

// NewPlayer 创建一个新玩家
func NewPlayer(id int64, nickname string, conn ziface.IConnection) *Player {
	return &Player{
		ID:       id,
		Nickname: nickname,
		Score:    1000, // 初始分数
		RoomID:   0,
		Hand:     make([]Card, 0),
		Status:   STATUS_WAITING,
		isOnline: true,
		Conn:     conn,
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

// SetOnline 设置玩家在线状态
func (p *Player) SetOnline(online bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isOnline = online
}

// IsOnline 检查玩家是否在线
func (p *Player) IsOnline() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isOnline
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

// SetBanker 设置玩家为庄家
func (p *Player) SetBanker(isBanker bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isBanker = isBanker
}

// IsBanker 检查玩家是否是庄家
func (p *Player) IsBanker() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isBanker
}

// PlaceBet 下注
func (p *Player) PlaceBet(amount int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.BetAmount = amount
}

// GetBetAmount 获取下注金额
func (p *Player) GetBetAmount() int32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.BetAmount
}

// HasBet 检查玩家是否已下注
func (p *Player) HasBet() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.BetAmount > 0
}

// ResetBet 重置下注
func (p *Player) ResetBet() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.BetAmount = 0
}

// PlayerStatus 玩家状态
type PlayerStatus int

const (
	STATUS_UNKNOWN PlayerStatus = iota
	STATUS_WAITING
	STATUS_READY
	STATUS_PLAYING
	STATUS_OFFLINE
)
