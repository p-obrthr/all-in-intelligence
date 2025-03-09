package tui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quit = true
			return m, tea.Quit
		}
		newView, cmd := m.view.HandleKey(msg)
		m.view = newView
		return m, cmd
	default:
		newView, cmd := m.view.HandleKey(msg)
		m.view = newView
		return m, cmd
	}
}
