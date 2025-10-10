package logic

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Card 表示一张扑克牌
type Card struct {
	Suit Suit
	Rank Rank
}

// Value 返回牌的点数，用于计算
func (c *Card) Value() int {
	if c.Rank >= RANK_TEN && c.Rank <= RANK_KING {
		return 10
	}
	return int(c.Rank)
}

// String 返回牌的字符串表示
func (c *Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank.String(), c.Suit.String())
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

func (s Suit) String() string {
	switch s {
	case SUIT_SPADES:
		return "Spades"
	case SUIT_HEARTS:
		return "Hearts"
	case SUIT_CLUBS:
		return "Clubs"
	case SUIT_DIAMONDS:
		return "Diamonds"
	default:
		return "Unknown"
	}
}

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

func (r Rank) String() string {
	switch r {
	case RANK_ACE:
		return "Ace"
	case RANK_TWO:
		return "Two"
	case RANK_THREE:
		return "Three"
	case RANK_FOUR:
		return "Four"
	case RANK_FIVE:
		return "Five"
	case RANK_SIX:
		return "Six"
	case RANK_SEVEN:
		return "Seven"
	case RANK_EIGHT:
		return "Eight"
	case RANK_NINE:
		return "Nine"
	case RANK_TEN:
		return "Ten"
	case RANK_JACK:
		return "Jack"
	case RANK_QUEEN:
		return "Queen"
	case RANK_KING:
		return "King"
	default:
		return "Unknown"
	}
}

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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.cards), func(i, j int) {
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
