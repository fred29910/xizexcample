package logic

import (
	"sort"
)

// CardType 牌型枚举
type CardType int

const (
	CARD_TYPE_NO_BULL       CardType = iota // 无牛
	CARD_TYPE_BULL_1                        // 牛一
	CARD_TYPE_BULL_2                        // 牛二
	CARD_TYPE_BULL_3                        // 牛三
	CARD_TYPE_BULL_4                        // 牛四
	CARD_TYPE_BULL_5                        // 牛五
	CARD_TYPE_BULL_6                        // 牛六
	CARD_TYPE_BULL_7                        // 牛七
	CARD_TYPE_BULL_8                        // 牛八
	CARD_TYPE_BULL_9                        // 牛九
	CARD_TYPE_BULL_BOMB                     // 牛牛 (炸弹)
	CARD_TYPE_FIVE_SMALL                    // 五小牛
	CARD_TYPE_BOMB                          // 炸弹 (四张相同)
	CARD_TYPE_GOLDEN_FLOWER                 // 金花 (同花顺)
)

// CalculateBull 计算牛牛牌型和牛值
func CalculateBull(cards []*Card) (CardType, uint32) {
	if len(cards) != 5 {
		return CARD_TYPE_NO_BULL, 0
	}

	// 检查特殊牌型
	if cardType, isSpecial := checkSpecialCardTypes(cards); isSpecial {
		return cardType, 0
	}

	// 检查是否有牛
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			for k := j + 1; k < 5; k++ {
				if (cards[i].Value()+cards[j].Value()+cards[k].Value())%10 == 0 {
					// 有牛，计算剩余两张牌的点数
					remainingSum := 0
					for idx, card := range cards {
						if idx != i && idx != j && idx != k {
							remainingSum += card.Value()
						}
					}
					bullValue := remainingSum % 10
					if bullValue == 0 {
						return CARD_TYPE_BULL_BOMB, 10 // 牛牛
					}
					return CardType(bullValue), uint32(bullValue)
				}
			}
		}
	}

	return CARD_TYPE_NO_BULL, 0
}

// checkSpecialCardTypes 检查特殊牌型
func checkSpecialCardTypes(cards []*Card) (CardType, bool) {
	// 检查五小牛: 五张牌点数之和小于等于10
	sum := 0
	for _, card := range cards {
		sum += card.Value()
	}
	if sum <= 10 {
		return CARD_TYPE_FIVE_SMALL, true
	}

	// 检查炸弹: 四张牌点数相同
	rankCount := make(map[Rank]int)
	for _, card := range cards {
		rankCount[card.Rank]++
	}
	for _, count := range rankCount {
		if count == 4 {
			return CARD_TYPE_BOMB, true
		}
	}

	// 检查金花: 同花顺
	if isFlush(cards) && isStraight(cards) {
		return CARD_TYPE_GOLDEN_FLOWER, true
	}

	// 检查同花
	if isFlush(cards) {
		// 同花但不是顺子，按牛值处理，这里简化为牛牛
		return CARD_TYPE_BULL_BOMB, true
	}

	// 检查顺子
	if isStraight(cards) {
		// 顺子但不是同花，按牛值处理，这里简化为牛牛
		return CARD_TYPE_BULL_BOMB, true
	}

	return CARD_TYPE_NO_BULL, false
}

// isFlush 判断是否是同花
func isFlush(cards []*Card) bool {
	suit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

// isStraight 判断是否是顺子
func isStraight(cards []*Card) bool {
	// 复制一份并排序
	sortedCards := make([]*Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return sortedCards[i].Rank < sortedCards[j].Rank
	})

	// 检查普通顺子
	for i := 1; i < len(sortedCards); i++ {
		if sortedCards[i].Rank != sortedCards[i-1].Rank+1 {
			// 检查 A-2-3-4-5 这种特殊情况
			if sortedCards[0].Rank == 1 && sortedCards[1].Rank == 2 && sortedCards[2].Rank == 3 && sortedCards[3].Rank == 4 && sortedCards[4].Rank == 5 {
				return true
			}
			// 检查 10-J-Q-K-A 这种特殊情况
			if sortedCards[0].Rank == 1 && sortedCards[1].Rank == 10 && sortedCards[2].Rank == 11 && sortedCards[3].Rank == 12 && sortedCards[4].Rank == 13 {
				return true
			}
			return false
		}
	}
	return true
}

// CompareHands 比较两手牌的大小，返回 1 表示 hand1 赢，-1 表示 hand2 赢，0 表示平局
func CompareHands(hand1, hand2 []*Card) int {
	type1, value1 := CalculateBull(hand1)
	type2, value2 := CalculateBull(hand2)

	// 先比较牌型
	if type1 > type2 {
		return 1
	} else if type1 < type2 {
		return -1
	}

	// 牌型相同，比较牛值
	if value1 > value2 {
		return 1
	} else if value1 < value2 {
		return -1
	}

	// 牛值也相同，比较最大单张牌
	maxCard1 := getMaxCard(hand1)
	maxCard2 := getMaxCard(hand2)

	if maxCard1.Value() > maxCard2.Value() {
		return 1
	} else if maxCard1.Value() < maxCard2.Value() {
		return -1
	}

	// 最大单张也相同，比较花色 (黑桃 > 红桃 > 梅花 > 方块)
	if maxCard1.Suit > maxCard2.Suit {
		return 1
	} else if maxCard1.Suit < maxCard2.Suit {
		return -1
	}

	return 0 // 完全相同
}

// getMaxCard 获取手牌中最大的牌
func getMaxCard(cards []*Card) *Card {
	maxCard := cards[0]
	for _, card := range cards {
		if card.Value() > maxCard.Value() {
			maxCard = card
		} else if card.Value() == maxCard.Value() && card.Suit > maxCard.Suit {
			maxCard = card
		}
	}
	return maxCard
}
