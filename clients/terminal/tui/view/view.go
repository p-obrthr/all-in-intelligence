package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	Init() tea.Cmd
	Render() string
	HandleKey(msg interface{}) (View, tea.Cmd)
}
