package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	game := newGame()

	doFlop(&game.Board, &game.Deck)
	doRiver(&game.Board, &game.Deck)
	doTurn(&game.Board, &game.Deck)
	p := tea.NewProgram(game)
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
