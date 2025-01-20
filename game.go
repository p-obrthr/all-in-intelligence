package main

import tea "github.com/charmbracelet/bubbletea"

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
		CurrentRound:   -1,
	}

	game.startNewRound()

	return game
}

func (round *Round) determineWinner() int {
	sortPlayersByRanking(&round.Players)
	for _, player := range round.Players {
		if !player.IsOut {
			return player.Id - 1
		}
	}
	return -1
}

func (game *Game) startNewRound() {
	if game.Quit || game.CurrentRound >= len(game.Rounds) {
		return
	}

	newRound := newRound(game.Players)
	game.Rounds = append(game.Rounds, newRound)
	game.CurrentRound++
}

func (game *Game) processGameUpdate(msg interface{}) {
	if game.WaitingForNextRound {
		if _, ok := msg.(tea.KeyMsg); ok {
			game.WaitingForNextRound = false
			game.startNewRound()
		}
	}

	round := &game.Rounds[game.CurrentRound]
	player := &round.Players[round.CurrentPlayer]
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		actionSuccessful, message := game.handleKeyInput(keyMsg)
		if actionSuccessful {
			round.MsgLog = append(round.MsgLog, colorMessage(actionSuccessful, message))

			activePlayers := getCountActivePlayers(*round)
			if activePlayers < 2 {
				game.endRound()
			}

			player.HasActed = true

			round.nextPlayer()

			if isAllPlayersHaveActed(*round) {
				if !round.nextStage() {
					game.endRound()
				}
			}
		} else {
			round.MsgLog = append(round.MsgLog, colorMessage(actionSuccessful, message))
		}
	}
}
