package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
	case " ":
		return game.Rounds[game.CurrentRound].check()
	case "q":
		game.Quit = true
		return false, ""
	default:
		return false, ""
	}
}

func (game Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if game.WaitingForNextRound {
		if _, ok := msg.(tea.KeyMsg); ok {
			game.WaitingForNextRound = false
			game.startNewRound()
		}
		return game, nil
	}

	if len(game.Rounds) == 0 || game.CurrentRound >= len(game.Rounds) {
		return game, nil
	}

	game.processGameUpdate(msg)

	return game, nil
}

func (game Game) View() string {
	if game.CurrentRound >= len(game.Rounds) {
		return "Invalid round."
	}

	round := &game.Rounds[game.CurrentRound]
	player := &round.Players[round.CurrentPlayer]

	checkStatus(player, &round.Board)

	boardDisplay := getCardsString(round.Board)
	playerCardsDisplay := getCardsString(player.Cards)

	return fmt.Sprintf(
		"Pot: %d\n\nBoard\n%s\n\nCurrent Player: %s (Money: %d)\n%s (%s)\n\nActions:\n[r] Raise\n[f] Fold\n[c] Call\n[_] Check\n[q] Quit\n\nLog:\n%s",
		round.Pot,
		boardDisplay,
		player.Name,
		player.Money,
		playerCardsDisplay,
		player.Status.TypeName,
		strings.Join(round.MsgLog, "\n"),
	)
}
