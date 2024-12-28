package main

import (
	"fmt"
)

type Game struct {
	Pot int
}

type Board [5]Card

type Deck []Card

func startGame() {
	deck := newDeck()

	var players []Player
	for i := 1; i <= 3; i++ {
		player := Player{
			Name: fmt.Sprintf("player%d", i),
		}
		drawPlayerCards(&player, &deck)
		players = append(players, player)
	}

	for _, player := range players {
		fmt.Printf("cards of %s:\n", player.Name)
		fmt.Println(player.Cards)
		fmt.Println("")
	}

	var board Board
	doFlop(&board, &deck)
	doTurn(&board, &deck)
	doRiver(&board, &deck)

	for i := range players {
		fmt.Printf("%s:\n", players[i].Name)
		checkStatus(&players[i], &board)
		fmt.Println("")
	}

	sortPlayersByRanking(&players)

	fmt.Println("final rankings:")
	for _, player := range players {
		fmt.Printf("%s (score: %d)\n", player.Name, player.Status.Score)
	}
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

func checkStatus(player *Player, board *Board) string {
	var cards []Card

	for i := 0; i < len(player.Cards); i++ {
		cards = append(cards, player.Cards[i])
	}

	for i := 0; i < len(board); i++ {
		cards = append(cards, board[i])
	}

	cards = filterEmptyCards(cards)

	// fmt.Println(cards)

	sortCardsByRanking(&cards)
	// fmt.Println("sorted:", cards)

	ranking, rankedCards := checkRanking(cards)

	fmt.Println(ranking)
	score := 0

	for _, card := range rankedCards {
		score += getCardRankValue(card)
	}
	player.Status = Status{
		Cards: rankedCards,
		Type:  getRankingType(ranking),
		Score: score,
	}

	fmt.Println(player.Status)

	return ""
}
