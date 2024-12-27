package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Card struct {
	Symbol string
	Rank   string
}

type Deck []Card

type Board [5]Card

type Player struct {
	Name   string
	Money  int
	Cards  [2]Card
	Status string
}

func main() {
	symbols := []string{"♣", "♠", "♥", "♦"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	// rankingCombinations := []string{
	// 	"highCard",
	// 	"pair",
	// 	"twoPair",
	// 	"threeOfAKind",
	// 	"straight",
	// 	"flush",
	// 	"fullHouse",
	// 	"fourOfAKind",
	// 	"straightFlush",
	// 	"royalFlush",
	// }

	var deck Deck

	for _, symbol := range symbols {
		for _, rank := range ranks {
			deck = append(deck, Card{Symbol: symbol, Rank: rank})
		}
	}

	// drawCard(&deck)
	// fmt.Println("total number:", len(deck))

	// drawCard(&deck)
	// fmt.Println("total number:", len(deck))

	// drawCard(&deck)
	// fmt.Println("total number:", len(deck))

	var board Board

	doFlop(&board, &deck)
	doTurn(&board, &deck)
	doRiver(&board, &deck)

	var playerOne Player
	playerOne.Name = "player one"
	drawPlayerCards(&playerOne, &deck)
	fmt.Println("cards of player one:")
	fmt.Println(playerOne.Cards[0])
	fmt.Println(playerOne.Cards[1])

	checkStatus(&playerOne, &board)
}

func drawCard(deck *Deck) Card {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(*deck))

	card := (*deck)[n]
	// fmt.Println("draw card:", card)

	*deck = append((*deck)[:n], (*deck)[n+1:]...)

	return card
}

func doFlop(board *Board, deck *Deck) {
	for i := 0; i < 3; i++ {
		board[i] = drawCard(deck)
	}
}

func doTurn(board *Board, deck *Deck) {
	board[3] = drawCard(deck)
}

func doRiver(board *Board, deck *Deck) {
	board[4] = drawCard(deck)
}

func drawPlayerCards(player *Player, deck *Deck) {
	for i := 0; i < 2; i++ {
		player.Cards[i] = drawCard(deck)
	}
}

func checkStatus(player *Player, board *Board) string {
	var cards []Card

	for i := 0; i < len(player.Cards); i++ {
		cards = append(cards, player.Cards[i])
	}

	for i := 0; i < len(board); i++ {
		cards = append(cards, board[i])
	}

	fmt.Println(cards)

	sortCardsByRanking(&cards)
	fmt.Println("sorted:", cards)

	fmt.Println(checkRanking(cards))

	return ""
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

	threeOfKinds := checkNOfAKind(cards, 3)

	if threeOfKinds == nil {
		return nil
	}

	twoPairs := checkNOfAKind(cards, 2)

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

	var sortedRanks []int
	for _, card := range cards {
		sortedRanks = append(sortedRanks, rankOrder[card.Rank])
	}

	rankSet := make(map[int]bool)
	for _, rank := range sortedRanks {
		rankSet[rank] = true
	}

	var uniqueRanks []int
	for rank := range rankSet {
		uniqueRanks = append(uniqueRanks, rank)
	}
	sort.Ints(uniqueRanks)

	for i := 0; i <= len(uniqueRanks)-5; i++ {
		if uniqueRanks[i+4]-uniqueRanks[i] == 4 {
			var straight []Card
			for _, rank := range uniqueRanks[i : i+5] {
				for _, card := range cards {
					if rankOrder[card.Rank] == rank {
						straight = append(straight, card)
						break
					}
				}
			}
			sortCardsByRanking(&straight)
			return straight
		}
	}

	lowStreet := []int{rankOrder["2"], rankOrder["3"], rankOrder["4"], rankOrder["5"]}
	if rankSet[rankOrder["A"]] {
		var straight []Card
		for _, rank := range lowStreet {
			for _, card := range cards {
				if rankOrder[card.Rank] == rank {
					straight = append(straight, card)
					break
				}
			}
		}
		if len(straight) == 5 {
			sortCardsByRanking(&straight)
			return straight
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
