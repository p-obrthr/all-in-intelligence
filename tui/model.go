package tui

import (
	"all-in-intelligence/tui/view"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	quit bool
	view view.View
}

func NewModel() *Model {
	return &Model{
		view: view.NewStartView(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.quit {
		return "Goodbye."
	}
	if m.view == nil {
		return "Loading..."
	}
	return "\n" + m.view.Render()
}
