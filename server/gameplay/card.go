package gameplay

import (
	"fmt"
	"math/rand"
	"sort"
)

type Card struct {
	Symbol string
	Rank   string
}

type Deck []Card

type Board [5]Card

func newDeck() Deck {
	symbols := []string{"♣", "♠", "♥", "♦"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	var deck Deck
	for _, symbol := range symbols {
		for _, rank := range ranks {
			deck = append(deck, Card{Symbol: symbol, Rank: rank})
		}
	}
	return deck
}

func drawCard(deck *Deck) Card {
	n := rand.Intn(len(*deck))

	card := (*deck)[n]
	*deck = append((*deck)[:n], (*deck)[n+1:]...)
	return card
}

func filterEmptyCards(cards []Card) []Card {
	var filtered []Card
	for _, card := range cards {
		if card.Symbol != "" && card.Rank != "" {
			filtered = append(filtered, card)
		}
	}
	return filtered
}

func checkRanking(cards []Card) (string, []Card) {
	var rankedCards []Card

	rankedCards = checkRoyalFlush(cards)
	if rankedCards != nil {
		return "royalFlush", rankedCards
	}

	rankedCards = checkStreetFlush(cards)
	if rankedCards != nil {
		return "straightFlush", rankedCards
	}

	rankedCards = checkNOfAKind(&cards, 4)
	if rankedCards != nil {
		return "fourOfAKind", rankedCards
	}

	rankedCards = checkFullHouse(&cards)
	if rankedCards != nil {
		return "fullHouse", rankedCards
	}

	rankedCards = checkFlush(cards)
	if rankedCards != nil {
		return "flush", rankedCards
	}

	rankedCards = checkStreet(cards)
	if rankedCards != nil {
		return "straight", rankedCards
	}

	rankedCards = checkNOfAKind(&cards, 3)
	if rankedCards != nil {
		return "threeOfAKind", rankedCards
	}

	rankedCards = checkNOfAKind(&cards, 2)
	if rankedCards != nil {
		if len(rankedCards) == 4 {
			return "twoPair", rankedCards
		}
		return "pair", rankedCards
	}

	return "highCard", cards[:1]
}

func checkRoyalFlush(cards []Card) []Card {
	flushs := checkFlush(cards)

	if flushs == nil {
		return nil
	}

	requiredRanks := map[string]bool{
		"10": false,
		"J":  false,
		"Q":  false,
		"K":  false,
		"A":  false,
	}

	for _, card := range flushs {
		if _, exists := requiredRanks[card.Rank]; exists {
			requiredRanks[card.Rank] = true
		}
	}

	if requiredRanks["10"] && requiredRanks["J"] && requiredRanks["Q"] && requiredRanks["K"] && requiredRanks["A"] {
		return flushs
	}
	return nil
}

func checkStreetFlush(cards []Card) []Card {
	flushs := checkFlush(cards)

	if flushs == nil {
		return nil
	}

	streets := checkStreet(flushs)

	if streets == nil {
		return nil
	} else {
		return streets
	}
}

func checkFullHouse(cards *[]Card) []Card {
	if len(*cards) < 5 {
		return nil
	}

	cardsCopy := append([]Card(nil), *cards...)

	threeOfKinds := checkNOfAKind(&cardsCopy, 3)

	if threeOfKinds == nil {
		return nil
	}

	twoPairs := checkNOfAKind(&cardsCopy, 2)

	if twoPairs == nil {
		return nil
	}

	var fullHouses []Card
	fullHouses = append(fullHouses, threeOfKinds...)
	fullHouses = append(fullHouses, twoPairs...)
	sortCardsByRanking(&fullHouses)
	return fullHouses
}

func checkFlush(cards []Card) []Card {
	if len(cards) < 5 {
		return nil
	}

	symbols := []string{"♣", "♠", "♥", "♦"}

	flushCards := make(map[string][]Card)

	for _, card := range cards {
		for _, symbol := range symbols {
			if card.Symbol == symbol {
				flushCards[symbol] = append(flushCards[symbol], card)
				break
			}
		}
	}

	for _, cardsSymbol := range flushCards {
		if len(cardsSymbol) >= 5 {
			return cardsSymbol
		}
	}

	return nil
}

func checkNOfAKind(cards *[]Card, n int) []Card {
	if len(*cards) < n {
		return nil
	}

	var nOfKinds []Card
	for i := 0; i <= len(*cards)-n; i++ {
		match := true

		for j := 0; j < n-1; j++ {
			if (*cards)[i+j].Rank != (*cards)[i+j+1].Rank {
				match = false
				break
			}
		}

		if match {
			nOfKinds = append(nOfKinds, (*cards)[i:i+n]...)
			*cards = append((*cards)[:i], (*cards)[i+n:]...)
			i--
			if len(nOfKinds) > 4 {
				// for two pair limited by four
				nOfKinds = nOfKinds[:4]
			}
		}
	}

	return nOfKinds
}

func checkStreet(cards []Card) []Card {
	if len(cards) < 5 {
		return nil
	}

	rankOrder := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10,
		"J": 11, "Q": 12, "K": 13, "A": 14,
	}

	rankSet := make(map[string]Card)
	for _, card := range cards {
		rankSet[card.Rank] = card
	}

	uniqueCards := []Card{}
	for _, card := range rankSet {
		uniqueCards = append(uniqueCards, card)
	}

	for i := 0; i <= len(uniqueCards)-5; i++ {
		isStreet := true
		streetCards := []Card{uniqueCards[i]}

		for j := 1; j < 5; j++ {
			if rankOrder[uniqueCards[i+j-1].Rank]-1 != rankOrder[uniqueCards[i+j].Rank] {
				isStreet = false
				break
			}
			streetCards = append(streetCards, uniqueCards[i+j])
		}

		if isStreet {
			return streetCards
		}
	}

	return nil
}

func sortCardsByRanking(cards *[]Card) {
	rankOrder := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10,
		"J": 11, "Q": 12, "K": 13, "A": 14,
	}
	sort.SliceStable(*cards, func(i, j int) bool {
		return rankOrder[(*cards)[i].Rank] > rankOrder[(*cards)[j].Rank]
	})
}

func getCardRankValue(card Card) int {
	rankOrder := map[string]int{
		"2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10,
		"J": 11, "Q": 12, "K": 13, "A": 14,
	}
	return rankOrder[card.Rank]
}

func getRankingType(ranking string) int {
	rankingMap := map[string]int{
		"highCard":      1,
		"pair":          2,
		"twoPair":       3,
		"threeOfAKind":  4,
		"straight":      5,
		"flush":         6,
		"fullHouse":     7,
		"fourOfAKind":   8,
		"straightFlush": 9,
		"royalFlush":    10,
	}
	return rankingMap[ranking]
}

func GetCards(cards []Card) string {
	cardsStr := ""
	for _, card := range cards {
		if card.Symbol != "" && card.Rank != "" {
			cardsStr += fmt.Sprintf("[%s %s] ", card.Rank, card.Symbol)
		} else {
			cardsStr += ""
		}
	}
	return cardsStr
}
