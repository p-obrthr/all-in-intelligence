package main

type Game struct {
	Players       []Player
	Board         Board
	Deck          Deck
	Pot           int
	LastRaise     int
	CurrentPlayer int
	MsgLog        []string
	Quit          bool
}

func newGame() Game {
	deck := newDeck()
	var players []Player

	totalPlayers := 3

	for i := 1; i <= totalPlayers; i++ {
		newPlayer := newPlayer(i)
		drawPlayerCards(&newPlayer, &deck)
		players = append(players, newPlayer)
	}

	return Game{
		Players:       players,
		Board:         Board{},
		Deck:          deck,
		Pot:           0,
		CurrentPlayer: 0,
		MsgLog:        []string{},
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
