package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (game Game) Init() tea.Cmd {
	return nil
}

func (game *Game) handleKeyInput(key tea.KeyMsg) (bool, string) {
	switch key.String() {
	case "r":
		return game.Rounds[game.CurrentRound].raise()
	case "f":
		return game.Rounds[game.CurrentRound].fold()
	case "c":
		return game.Rounds[game.CurrentRound].call()
	case " ", "_":
		return game.Rounds[game.CurrentRound].check()
	default:
		return false, ""
	}
}

func (game *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch game.Phase {
	case PhaseConfiguring:
		return game.updateConfig(msg)
	case PhasePlaying:
		return game.updatePlaying(msg)
	}
	return game, nil
}

func (game *Game) updateConfig(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &game.Config.questions[game.Config.index]

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return game, tea.Quit
		case "enter":
			current.answer = game.Config.answerField.Value()
			game.NextQuestion()

			if game.Config.done {
				game.InitGame()
				game.Phase = PhasePlaying
			}

			return game, nil
		}
	}

	var cmd tea.Cmd
	game.Config.answerField, cmd = game.Config.answerField.Update(msg)
	return game, cmd
}

func (game *Game) updatePlaying(msg tea.Msg) (tea.Model, tea.Cmd) {
	if game.WaitingForNextRound {
		if _, ok := msg.(tea.KeyMsg); ok {
			game.WaitingForNextRound = false
			game.startNewRound()
			return game, tea.Batch(tea.EnterAltScreen)
		}
		return game, nil
	}

	round := &game.Rounds[game.CurrentRound]
	player := round.Players[round.CurrentPlayer]

	if player.IsLLM {
		action := GetLLMAction(*round)
		game.applyLLMAction(action)
		return game, nil
	}

	cmd := game.processGameUpdate(msg)
	return game, cmd
}

func (game *Game) NextQuestion() {
	game.Config.questions[game.Config.index].answer = game.Config.answerField.Value()

	if game.Config.index < len(game.Config.questions)-1 {
		game.Config.index++
		game.Config.answerField = textinput.New()
		game.Config.answerField.Focus()
	} else {
		game.Config.done = true
	}
}

func (game *Game) applyLLMAction(action string) {
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(action)}
	actionSuccessful, message := game.handleKeyInput(keyMsg)
	game.processPlayerAction(actionSuccessful, message)
}

func (game Game) View() string {
	if game.Phase == PhaseConfiguring {

		cardArt := `
       _____      _____      _____      _____
      |A .  |    |K .  |    |Q .  |    |J .  |
      | /.\ ||   | /.\ ||   | /.\ ||   | /.\ ||
      |(_._)||   |(_._)||   |(_._)||   |(_._)||
      |  |  ||   |  |  ||   |  |  ||   |  |  ||
      |____V||   |____V||   |____V||   |____V||
             |____V||   |____V||   |____V||   
             |____V||   |____V||   |____V||   
      .__  .__    .__         .__        __         .__  .__  .__                                   
    _____  |  | |  |   |__| ____   |__| _____/  |_  ____ |  | |  | |  | |__| ____   ____   ____   ____  
    \__  \ |  | |  |   |  |/    \  |  |/    \   __\/ __ \|  | |  | |  |/ ___\_/ __ \ /    \_/ ___\/ __ \ 
     / __ \|  |_|  |__ |  |   |  \ |  |   |  \  | \  ___/|  |_|  |_|  / /_/  >  ___/|   |  \  \__\  ___/ 
    (____  /____/____/ |__|___|  / |__|___|  /__|  \___  >____/____/__\___  / \___  >___|  /\___  >___  >
         \/                    \/          \/          \/            /_____/      \/     \/     \/    \/  
	`

		return lipgloss.JoinVertical(
			lipgloss.Center,
			cardArt,
			game.Config.questions[game.Config.index].question,
			game.Config.styles.InputField.Render(game.Config.answerField.View()),
		)
	}

	if game.CurrentRound >= len(game.Rounds) {
		return "Invalid round."
	}

	round := &game.Rounds[game.CurrentRound]
	player := &round.Players[round.CurrentPlayer]

	checkStatus(player, &round.Board)

	boardDisplay := getCardsString(round.Board)
	playerCardsDisplay := getCardsString(player.Cards)

	return fmt.Sprintf(
		"Pot: %d\n\nBoard\n%s\n\nCurrent Player: %s (Money: %d)\n%s (%s)\n\nActions:\n[r] Raise\n[f] Fold\n[c] Call\n[_] Check\n\nLog:\n%s",
		round.Pot,
		boardDisplay,
		player.Name,
		player.Money,
		playerCardsDisplay,
		player.Status.TypeName,
		strings.Join(round.MsgLog, "\n"),
	)
}
