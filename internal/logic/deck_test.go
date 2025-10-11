package logic

import (
	"reflect"
	"testing"
)

func TestDeckNewDeck(t *testing.T) {
	deck := NewDeck()

	if len(deck.cards) != 52 {
		t.Errorf("Expected new deck to have 52 cards, got %d", len(deck.cards))
	}

	// 检查是否每种花色和点数都有
	cardCounts := make(map[string]int)
	for _, card := range deck.cards {
		cardCounts[card.String()]++
	}

	if len(cardCounts) != 52 {
		t.Errorf("Expected 52 unique cards, got %d", len(cardCounts))
	}
}

func TestDeckShuffle(t *testing.T) {
	deck1 := NewDeck()
	deck2 := NewDeck()

	// 确保两个未洗牌的牌堆顺序相同
	if !reflect.DeepEqual(deck1.cards, deck2.cards) {
		t.Error("Two new decks should have the same order before shuffling")
	}

	// 洗牌
	deck1.Shuffle()

	// 洗牌后顺序应该不同
	if reflect.DeepEqual(deck1.cards, deck2.cards) {
		t.Error("Shuffled deck should have a different order than an unshuffled one")
	}

	// 洗牌后牌的数量应该不变
	if len(deck1.cards) != 52 {
		t.Errorf("Expected shuffled deck to still have 52 cards, got %d", len(deck1.cards))
	}
}

func TestDeckDeal(t *testing.T) {
	deck := NewDeck()

	// 测试发牌
	cards, err := deck.DealCards(5)
	if err != nil {
		t.Fatalf("Deal failed: %v", err)
	}

	if len(cards) != 5 {
		t.Errorf("Expected to deal 5 cards, got %d", len(cards))
	}

	// 检查牌堆中的牌是否减少
	if len(deck.cards) != 52-5 {
		t.Errorf("Expected deck to have 47 cards after dealing 5, got %d", len(deck.cards))
	}

	// 测试从空牌堆发牌
	for i := 0; i < 47; i++ {
		_, err := deck.DealCard()
		if err != nil {
			t.Fatalf("DealCard failed during emptying deck: %v", err)
		}
	}

	// 尝试再发一张牌，应该失败
	_, err = deck.DealCard()
	if err == nil {
		t.Error("Expected error when dealing from an empty deck, but got nil")
	}

	// 尝试发比剩余牌更多的牌
	deck.Reset()
	_, err = deck.DealCards(53)
	if err == nil {
		t.Error("Expected error when dealing more cards than available, but got nil")
	}
}

func TestDeckReset(t *testing.T) {
	deck := NewDeck()

	// 发一些牌
	deck.DealCards(10)

	if len(deck.cards) != 42 {
		t.Errorf("Expected deck to have 42 cards after dealing 10, got %d", len(deck.cards))
	}

	// 重置牌堆
	deck.Reset()

	if len(deck.cards) != 52 {
		t.Errorf("Expected deck to have 52 cards after reset, got %d", len(deck.cards))
	}

	// 检查是否是新的完整牌堆（顺序可能不同）
	cardCounts := make(map[string]int)
	for _, card := range deck.cards {
		cardCounts[card.String()]++
	}

	if len(cardCounts) != 52 {
		t.Errorf("Expected 52 unique cards after reset, got %d", len(cardCounts))
	}
}
