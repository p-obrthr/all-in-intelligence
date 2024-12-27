package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Card struct {
	Symbol string
	Rank   string
}

type Deck []Card

type Board [5]Card

type Player struct {
	Name  string
	Money int
	Cards [2]Card
}

func main() {
	symbols := []string{"♣", "♠", "♥", "♦"}
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	var deck Deck

	for _, suit := range symbols {
		for _, rank := range ranks {
			deck = append(deck, Card{Symbol: suit, Rank: rank})
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
}

func drawCard(deck *Deck) Card {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(*deck))

	card := (*deck)[n]
	fmt.Println("draw card:", card)

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
