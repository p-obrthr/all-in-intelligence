package main

import (
	"fmt"
	"os"

	"terminal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := tui.NewModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("err: running program:", err)
		os.Exit(1)
	}
}
