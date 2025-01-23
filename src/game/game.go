package game

import (
	"src/services"

	tea "github.com/charmbracelet/bubbletea"
)

type Game struct {
	Players             []Player
	Rounds              []Round
	CurrentRound        int
	Quit                bool
	BigBlindAmount      int
	WaitingForNextRound bool
	OpenAIClient        services.OpenAIClient
}

func NewGame(client services.OpenAIClient) Game {
	var players []Player

	totalPlayers := 2

	for i := 1; i <= totalPlayers; i++ {
		newPlayer := newPlayer(i)
		if i%2 == 0 {
			newPlayer.IsLLM = true
			newPlayer.Name += " (LLM)"
		}
		players = append(players, newPlayer)
	}

	game := Game{
		Players:        players,
		Rounds:         []Round{},
		BigBlindAmount: 40,
		Quit:           false,
		CurrentRound:   -1,
		OpenAIClient:   client,
	}

	game.startNewRound()

	return game
}

func (round *Round) determineWinner() *Player {
	sortPlayersByRanking(&round.Players)
	for _, player := range round.Players {
		if !player.IsOut {
			return &player
		}
	}
	return nil
}

func (game *Game) startNewRound() {
	if game.Quit || game.CurrentRound >= len(game.Rounds) {
		return
	}

	newRound := newRound(game.Players, game.OpenAIClient)
	game.Rounds = append(game.Rounds, newRound)
	game.CurrentRound++
}

func (game *Game) processGameUpdate(msg interface{}) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		actionSuccessful, message := game.handleKeyInput(keyMsg)
		game.processPlayerAction(actionSuccessful, message)
		return tea.Batch(tea.EnterAltScreen)
	}
	return nil
}
