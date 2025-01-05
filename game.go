package main

type Game struct {
	Players             []Player
	Rounds              []Round
	CurrentRound        int
	Quit                bool
	BigBlindAmount      int
	WaitingForNextRound bool
}

func newGame() Game {
	var players []Player

	totalPlayers := 3

	for i := 1; i <= totalPlayers; i++ {
		newPlayer := newPlayer(i)
		players = append(players, newPlayer)
	}

	game := Game{
		Players:        players,
		Rounds:         []Round{},
		BigBlindAmount: 40,
		Quit:           false,
		CurrentRound:   0,
	}

	newRound := newRound(game.Players)
	game.Rounds = append(game.Rounds, newRound)

	return game
}

func (round *Round) determineWinner() int {
	sortPlayersByRanking(&round.Players)

	for i := 0; i < len(round.Players); i++ {
		if !round.Players[i].IsOut {
			return round.Players[i].Id - 1
		}
	}
	return -1
}

func (game *Game) startNewRound() {
	if game.CurrentRound >= len(game.Rounds) {
		game.Quit = true
		return
	}

	newRound := newRound(game.Players)
	game.Rounds = append(game.Rounds, newRound)
	game.CurrentRound++
}
