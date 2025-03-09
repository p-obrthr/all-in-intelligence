package main

import (
	"fmt"
	"os"

	"all-in-intelligence/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.NewModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
