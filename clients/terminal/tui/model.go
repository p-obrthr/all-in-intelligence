package tui

import (
	"encoding/json"
	"fmt"
	"log"
	"terminal/tui/view"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

type Model struct {
	quit     bool
	view     view.View
	wsConn   *websocket.Conn
	wsURL    string
	messages []string
}

type NewMessage struct {
	Message string
}

func NewModel() *Model {
	wsURL := "ws://localhost:8080/ws"

	model := &Model{
		view:  view.NewStartView(),
		wsURL: wsURL,
	}

	return model
}

func (m *Model) Init() tea.Cmd {
	cmd := m.view.Init()
	wsCmd := m.startWebSocketConnection()
	return tea.Batch(cmd, wsCmd)
}

func (m Model) View() string {
	if m.quit {
		if m.wsConn != nil {
			m.wsConn.Close()
		}
		return ""
	}
	if m.view == nil {
		return ""
	}

	return "\n" + m.view.Render()
}

func (m *Model) Close() {
	if m.wsConn != nil {
		m.wsConn.Close()
	}
}

func (m *Model) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		_, msg, err := m.wsConn.ReadMessage()
		if err != nil {
			fmt.Println("err receiving msg:", err)
			return nil
		}
		return NewMessage{Message: string(msg)}
	}
}

func (m *Model) startWebSocketConnection() tea.Cmd {
	conn, _, err := websocket.DefaultDialer.Dial(m.wsURL, nil)
	if err != nil {
		log.Fatal("err start ws connection:", err)
		return nil
	}
	m.wsConn = conn

	return m.listenForMessages()
}

func (m *Model) SendToWs(msg interface{}) error {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("err: marsh message:", err)
		return err
	}

	if m.wsConn == nil {
		return fmt.Errorf("err: ws connection not available")
	}

	err = m.wsConn.WriteMessage(websocket.TextMessage, msgJSON)
	if err != nil {
		fmt.Println("err: sending message to websocket:", err)
		return err
	}

	return nil
}
