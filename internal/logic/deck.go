package logic

import (
	"math/rand"
	"time"
)

// Deck 表示一副扑克牌
type Deck struct {
	cards []Card
}

// NewDeck 创建一副新的扑克牌
func NewDeck() *Deck {
	d := &Deck{
		cards: make([]Card, 0, 52),
	}
	d.Reset()
	return d
}

// Reset 重置并填充一副新牌
func (d *Deck) Reset() {
	d.cards = d.cards[:0] // 清空切片
	for suit := SUIT_SPADES; suit <= SUIT_DIAMONDS; suit++ {
		for rank := RANK_ACE; rank <= RANK_KING; rank++ {
			d.cards = append(d.cards, Card{Suit: suit, Rank: rank})
		}
	}
}

// Shuffle 洗牌
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// DealCard 发一张牌
func (d *Deck) DealCard() (Card, error) {
	if len(d.cards) == 0 {
		return Card{}, errors.New("deck is empty")
	}
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card, nil
}

// DealCards 发指定数量的牌
func (d *Deck) DealCards(num int) ([]Card, error) {
	if num > len(d.cards) {
		return nil, errors.New("not enough cards in deck")
	}
	cards := make([]Card, num)
	for i := 0; i < num; i++ {
		card, err := d.DealCard()
		if err != nil {
			return nil, err
		}
		cards[i] = card
	}
	return cards, nil
}

// GetCardCount 获取牌堆中剩余的牌数
func (d *Deck) GetCardCount() int {
	return len(d.cards)
}