package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (game Game) Init() tea.Cmd {
	return nil
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

	round := &game.Rounds[game.CurrentRound]
	player := &round.Players[round.CurrentPlayer]

	actionSuccessful := false
	message := ""
	switch msg := msg.(type) {
	case tea.KeyMsg:
		game.WaitingForNextRound = false
		switch msg.String() {
		case "r":
			actionSuccessful, message = round.raise()
		case "f":
			actionSuccessful, message = round.fold()
		case "c":
			actionSuccessful, message = round.call()
		case " ":
			actionSuccessful, message = round.check()
		case "q":
			game.Quit = true
			return game, tea.Quit
		}

		if actionSuccessful {
			round.MsgLog = append(round.MsgLog, Green.Render(message))

			activePlayers := getCountActivePlayers(*round)
			if activePlayers < 2 {
				game.endRound()
				return game, nil
			}

			player.HasActed = true

			totalPlayers := len(round.Players)
			for i := 1; i < totalPlayers; i++ {
				nextPlayerId := (round.CurrentPlayer + i) % totalPlayers
				nextPlayer := round.Players[nextPlayerId]
				if !nextPlayer.IsOut {
					round.CurrentPlayer = nextPlayerId
					break
				}
			}

			if allPlayersHaveActed(*round) {
				if !round.nextStage() {
					game.endRound()
				}
			}
		} else {
			round.MsgLog = append(round.MsgLog, Red.Render(message))
		}
	}

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
