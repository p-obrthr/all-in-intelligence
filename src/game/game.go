package game

import (
	"src/services"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Game struct {
	Config              Config
	Players             []Player
	Rounds              []Round
	CurrentRound        int
	Quit                bool
	BigBlindAmount      int
	StartingMoney       int
	WaitingForNextRound bool
	OpenAIClient        services.OpenAIClient
	Phase               Phase
}

type Phase int

const (
	PhaseConfiguring Phase = iota
	PhasePlaying
)

type Config struct {
	index       int
	questions   []Question
	answerField textinput.Model
	styles      *Styles
	done        bool
}

func NewConfig(questions []Question) *Config {
	styles := DefaultConfigStyles()
	answerField := textinput.New()
	answerField.Focus()
	return &Config{
		questions:   questions,
		answerField: answerField,
		styles:      styles,
	}
}

type Question struct {
	question string
	answer   string
}

func NewQuestion(question string) Question {
	return Question{question: question}
}

func New(client services.OpenAIClient) Game {
	questions := []Question{
		NewQuestion("How many players?"),
		NewQuestion("Specify which players are LLM Players (enter player numbers separated by commas):"),
		NewQuestion("Starting money for each player?"),
		NewQuestion("Big blind amount?"),
	}

	var game Game
	game.Config = *NewConfig(questions)
	game.OpenAIClient = client
	game.Phase = PhaseConfiguring

	return game
}

func (game *Game) InitGame() {
	totalPlayers, _ := strconv.Atoi(game.Config.questions[0].answer)

	strLlmPlayers := strings.Split(game.Config.questions[1].answer, ",")
	var llmPlayers []int
	for _, str := range strLlmPlayers {
		player, err := strconv.Atoi(str)
		if err != nil {
			continue
		}
		llmPlayers = append(llmPlayers, player)
	}

	startingMoney, _ := strconv.Atoi(game.Config.questions[2].answer)
	bigBlind, _ := strconv.Atoi(game.Config.questions[3].answer)

	game.BigBlindAmount = bigBlind

	for i := 1; i <= totalPlayers; i++ {
		newPlayer := newPlayer(
			i,
			startingMoney,
			contains(llmPlayers, i),
		)
		game.Players = append(game.Players, newPlayer)
	}

	game.CurrentRound = -1
	game.startNewRound()
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
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

	newRound := newRound(game.Players, game.OpenAIClient, game.BigBlindAmount)
	game.Rounds = append(game.Rounds, newRound)
	game.CurrentRound++
}

func (game *Game) processGameUpdate(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC {
			return tea.Quit
		}
		actionSuccessful, message := game.handleKeyInput(keyMsg)
		game.processPlayerAction(actionSuccessful, message)
		return tea.Batch(tea.EnterAltScreen)
	}
	return nil
}
