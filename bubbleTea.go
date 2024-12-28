package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	Red    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Grey   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	Green  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	Black  = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
)

func (game Game) Init() tea.Cmd {
	return nil
}

func (game Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		actionSuccessful := false
		switch msg.String() {
		case "r":
			var raiseAmount int
			if game.Pot*2 == 0 {
				raiseAmount = 10
			} else {
				raiseAmount = game.Pot * 2
			}
			if game.doRaise(raiseAmount) {
				game.MsgLog = append(game.MsgLog, Green.Render(fmt.Sprintf("%s raised to %d.", game.Players[game.CurrentPlayer].Name, raiseAmount)))
				actionSuccessful = true
			} else {
				game.MsgLog = append(game.MsgLog, Red.Render("not enough"))
			}
		case "f":
			if game.doFold() {
				game.MsgLog = append(game.MsgLog, Black.Render(fmt.Sprintf("%s folded.", game.Players[game.CurrentPlayer].Name)))
				actionSuccessful = true
			}
		case "c":
			if game.doCall() {
				game.MsgLog = append(game.MsgLog, Yellow.Render(fmt.Sprintf("%s called to %d.", game.Players[game.CurrentPlayer].Name, game.Players[game.CurrentPlayer].InPot)))
				actionSuccessful = true
			}
		case " ":
			game.MsgLog = append(game.MsgLog, Grey.Render(fmt.Sprintf("%s checked.", game.Players[game.CurrentPlayer].Name)))
			actionSuccessful = true
		case "q":
			game.Quit = true
			return game, tea.Quit
		}

		if actionSuccessful {
			totalPlayers := len(game.Players)
			for i := 1; i < totalPlayers; i++ {
				nextPlayer := (game.CurrentPlayer + i) % totalPlayers
				if !game.Players[nextPlayer].IsOut {
					game.CurrentPlayer = nextPlayer
					break
				}
			}
		}
	}
	return game, nil
}

func (game Game) View() string {
	if game.Quit {
		return "Goodbye."
	}

	current := game.Players[game.CurrentPlayer]

	checkStatus(&current, &game.Board)

	boardDisplay := getCardsString(game.Board[:])
	playerCardsDisplay := getCardsString(current.Cards[:])

	return fmt.Sprintf(
		"Pot: %d\n\nBoard\n%s\n\nCurrent Player: %s (Money: %d)\n%s (%s)\n\nActions:\n[r] Raise\n[f] Fold\n[c] Call\n[_] Check\n[q] Quit\n\nLog:\n%s",
		game.Pot,
		boardDisplay,
		current.Name,
		current.Money,
		playerCardsDisplay,
		current.Status.TypeName,
		strings.Join(game.MsgLog, "\n"),
	)
}
