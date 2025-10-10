package logic

import (
	"errors"
	"sync"
)

// Room 表示一个游戏房间
type Room struct {
	ID      int32
	Players map[int64]*Player // key: playerID
	Deck    *Deck
	mu      sync.RWMutex
}

// NewRoom 创建一个新房间
func NewRoom(roomID int32) *Room {
	return &Room{
		ID:      roomID,
		Players: make(map[int64]*Player),
		Deck:    NewDeck(),
	}
}

// AddPlayer 添加一个玩家到房间
func (r *Room) AddPlayer(player *Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Players) >= 5 {
		return errors.New("room is full")
	}

	if _, exists := r.Players[player.ID]; exists {
		return errors.New("player already in room")
	}

	player.SetRoomID(r.ID)
	r.Players[player.ID] = player
	return nil
}

// RemovePlayer 从房间移除一个玩家
func (r *Room) RemovePlayer(playerID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Players[playerID]; !exists {
		return errors.New("player not in room")
	}

	player := r.Players[playerID]
	player.SetRoomID(0) // 从房间中移除
	delete(r.Players, playerID)
	return nil
}

// GetPlayer 获取房间内的一个玩家
func (r *Room) GetPlayer(playerID int64) (*Player, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	player, exists := r.Players[playerID]
	if !exists {
		return nil, errors.New("player not found in room")
	}
	return player, nil
}

// GetPlayers 获取房间内所有玩家
func (r *Room) GetPlayers() []*Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	players := make([]*Player, 0, len(r.Players))
	for _, p := range r.Players {
		players = append(players, p)
	}
	return players
}

// GetPlayerCount 获取房间内玩家数量
func (r *Room) GetPlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Players)
}

// IsFull 检查房间是否已满
func (r *Room) IsFull() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.Players) >= 5
}

// ResetDeck 重置并洗牌
func (r *Room) ResetDeck() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Deck.Reset()
	r.Deck.Shuffle()
}

// DealCardsToPlayer 给指定玩家发牌
func (r *Room) DealCardsToPlayer(playerID int64, num int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	player, exists := r.Players[playerID]
	if !exists {
		return errors.New("player not found in room")
	}

	cards, err := r.Deck.DealCards(num)
	if err != nil {
		return err
	}

	player.ClearHand()
	for _, card := range cards {
		player.AddCard(card)
	}
	return nil
}