package view

import (
	"all-in-intelligence/gameplay"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type GameView struct {
	game gameplay.Game
}

func NewGameView(players int, bbAmount int) *GameView {
	v := GameView{
		game: gameplay.NewGame(players, bbAmount),
	}
	return &v
}

func (v GameView) Render() string {
	if v.game.CurrentRound >= len(v.game.Rounds) {
		return "Invalid round."
	}

	round := &v.game.Rounds[v.game.CurrentRound]
	player := &round.Players[round.CurrentPlayer]

	// game.checkStatus(player, &round.Board)

	boardDisplay := gameplay.GetCardsString(round.Board)
	playerCardsDisplay := gameplay.GetCardsString(player.Cards)

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

func (v *GameView) HandleKey(msg interface{}) (View, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		round := &v.game.Rounds[v.game.CurrentRound]

		var actionSuccessful bool
		var message string

		switch keyMsg.String() {
		case "r":
			actionSuccessful, message = round.Raise()
		case "f":
			actionSuccessful, message = round.Fold()
		case "c":
			actionSuccessful, message = round.Call()
		case " ", "_":
			actionSuccessful, message = round.Check()
		default:
			return v, nil
		}

		v.game.ProcessPlayerAction(actionSuccessful, message)
	}
	return v, nil
}
