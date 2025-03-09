package gameplay

import (
	"all-in-intelligence/services"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Game struct {
	Players             []Player
	Rounds              []Round
	CurrentRound        int
	Quit                bool
	BigBlindAmount      int
	StartingMoney       int
	WaitingForNextRound bool
	OpenAIClient        services.OpenAIClient
}

func NewGame(totalPlayers int, bbAmount int) Game {

	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		apiKey = ""
	}
	systemPrompt := "You are a professional poker player. You will be given a situation, and you must decide your action. Respond only in JSON format using the following options: { 'action': 'r' } for Raise, { 'action': 'f' } for Fold, { 'action': 'c' } for Call, and { 'action': '_' } for Check."
	client := services.NewOpenAIClient(apiKey, systemPrompt)

	game := Game{
		BigBlindAmount: bbAmount,
		OpenAIClient:   *client,
		CurrentRound:   -1,
	}

	startingMoney := 4000
	game.StartingMoney = startingMoney

	newMainPlayer := newPlayer(
		1,
		startingMoney,
		false,
	)
	game.Players = append(game.Players, newMainPlayer)

	for i := 2; i <= totalPlayers; i++ {
		newLlmPlayer := newPlayer(
			i,
			startingMoney,
			true,
		)
		game.Players = append(game.Players, newLlmPlayer)
	}

	game.startNewRound()
	return game
}

func (game *Game) startNewRound() {
	if game.Quit || game.CurrentRound >= len(game.Rounds) {
		return
	}

	newRound := newRound(game.Players, game.OpenAIClient, game.BigBlindAmount)
	game.Rounds = append(game.Rounds, newRound)
	game.CurrentRound++
}

func (game *Game) processGameUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC {
			return tea.Quit
		}
	}
	return nil
}

func (game *Game) updatePlaying(msg tea.Msg) {
	if game.WaitingForNextRound {
		if _, ok := msg.(tea.KeyMsg); ok {
			game.WaitingForNextRound = false
			game.startNewRound()
		}
	}

	round := &game.Rounds[game.CurrentRound]
	player := round.Players[round.CurrentPlayer]

	if player.IsLLM {
		action := GetLLMAction(*round)
		game.applyLLMAction(action)
	}

	// cmd := game.processGameUpdate(msg)
}

func (game *Game) applyLLMAction(action string) {
	// keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(action)}
	// actionSuccessful, message := game.handleKeyInput(keyMsg)
	// game.ProcessPlayerAction(actionSuccessful, message)
}
