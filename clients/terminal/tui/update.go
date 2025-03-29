package tui

import (
	"log"
	"terminal/tui/message"
	"terminal/tui/view"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case message.Config:
		// fmt.Println("case message sendConfigMsg")
		err := m.SendToWs(msg)
		if err != nil {
			log.Println("err sending config:", err)
		}
		return m, nil

	case message.Move:
		err := m.SendToWs(msg)
		if err != nil {
			log.Println("err sending move:", err)
		}
		return m, nil

	case NewMessage:
		if startView, ok := m.view.(*view.StartView); ok {
			startView.UpdateMessage(string(msg.Message))
		}
		if gameView, ok := m.view.(*view.GameView); ok {
			gameView.UpdateMessage(string(msg.Message))
		}
		return m, m.listenForMessages()

	default:
		// m.startWebSocketConnection()
		newView, cmd := m.view.HandleKey(msg)
		m.view = newView
		return m, cmd
	}
}
