package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	Render() string
	HandleKey(msg interface{}) (View, tea.Cmd)
}
