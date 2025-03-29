package view

import (
	"encoding/json"
	"fmt"
	"strings"
	"terminal/tui/message"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	Red       = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Grey      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	Green     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	Black     = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
)

type GameView struct {
	state  State
	keymap keymap
	help   help.Model
}

type State struct {
	Type            string   `json:"type"`
	Pot             int      `json:"pot"`
	Board           []Card   `json:"board"`
	Cards           []Card   `json:"cards"`
	StatusName      string   `json:"status_name"`
	PlayerId        int      `json:"player_id"`
	CurrentPlayerId int      `json:"current_player_id"`
	Money           int      `json:"money"`
	MsgLog          []string `json:"msg_log"`
}

type Card struct {
	Symbol string
	Rank   string
}

type Board [5]Card

type keymap struct {
	check key.Binding
	call  key.Binding
	fold  key.Binding
	raise key.Binding
}

func NewGameView() *GameView {
	return &GameView{
		keymap: keymap{

			check: key.NewBinding(
				key.WithKeys(" ", "_"),
				key.WithHelp("_", "check"),
			),
			call: key.NewBinding(
				key.WithKeys("c"),
				key.WithHelp("c", "call"),
			),
			fold: key.NewBinding(
				key.WithKeys("f"),
				key.WithHelp("f", "fold"),
			),
			raise: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "raise"),
			),
		},
		help: help.New(),
	}
}

func (v *GameView) Init() tea.Cmd {
	return nil
}

func (v *GameView) UpdateMessage(msg string) {
	var state State
	err := json.Unmarshal([]byte(msg), &state)
	if err != nil {
		logMessage := fmt.Sprintf("Fehler beim Unmarshallen der Nachricht: %v. Nachricht: %s", err, msg)
		fmt.Println(logMessage)
		return
	}

	v.state = state
	v.Render()
}

func (v GameView) Render() string {
	return "\n" + v.GetGameView()
}

func (v GameView) GetGameView() string {
	boardDisplay := GetCardsString(v.state.Board)

	playerCardsDisplay := ""
	if len(v.state.Cards) > 0 {
		playerCardsDisplay = GetCardsString(v.state.Cards)
	}

	currentPlayerDisplay := fmt.Sprintf("%d", v.state.CurrentPlayerId)
	if v.state.CurrentPlayerId == v.state.PlayerId && v.state.PlayerId != 0 {
		currentPlayerDisplay = Green.Render(fmt.Sprintf("%d (YOUR TURN)", v.state.CurrentPlayerId))
	}

	statusNameDisplay := v.state.StatusName
	if statusNameDisplay != "" {
		statusNameDisplay = fmt.Sprintf("(%s)", statusNameDisplay)
	}

	view := fmt.Sprintf(
		"Your Id: %d\n\n"+
			"Pot: %d\n\n"+
			"Board:\n%s\n\n"+
			"Your Cards:\n%s%s\n\n"+
			"Current Player Id: %s\n\n"+
			"Current Player Money: %d\n\n"+
			"Log:\n%s\n\n"+
			"%s",
		v.state.PlayerId,
		v.state.Pot,
		boardDisplay,
		playerCardsDisplay,
		statusNameDisplay,
		currentPlayerDisplay,
		v.state.Money,
		strings.Join(v.state.MsgLog, "\n"),
		v.helpView(),
	)

	return view
}

func (v GameView) helpView() string {
	return "\n" + v.help.ShortHelpView([]key.Binding{
		v.keymap.check,
		v.keymap.call,
		v.keymap.fold,
		v.keymap.raise,
	})
}

func GetCardString(card Card) string {
	var color lipgloss.Style
	switch card.Symbol {
	case "♠", "♣":
		color = Black
	case "♥", "♦":
		color = Red
	default:
		color = Grey
	}

	cardStr := "┌───────┐\n"
	cardStr += fmt.Sprintf("│%2s     │\n", card.Rank)
	cardStr += fmt.Sprintf("│   %s   │\n", color.Render(card.Symbol))
	cardStr += fmt.Sprintf("│     %2s│\n", card.Rank)
	cardStr += "└───────┘\n"
	cardStr += "\n"
	return cardStr
}

func GetCardsString(cards []Card) string {
	if len(cards) == 0 {
		return ""
	}

	var one, two, three, four, five string

	for _, card := range cards {
		cardStr := GetCardString(card)
		lines := strings.Split(cardStr, "\n")

		one += lines[0] + " "
		two += lines[1] + " "
		three += lines[2] + " "
		four += lines[3] + " "
		five += lines[4] + " "
	}

	result := strings.TrimSpace(one) + "\n" + strings.TrimSpace(two) + "\n" + strings.TrimSpace(three) + "\n" + strings.TrimSpace(four) + "\n" + strings.TrimSpace(five) + "\n"

	return result
}

func (v *GameView) HandleKey(msg interface{}) (View, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		key := keyMsg.String()

		validActions := map[string]string{
			"r": "r",
			"f": "f",
			"c": "c",
			" ": "_",
			"_": "_",
		}

		if action, found := validActions[key]; found {
			move := message.Move{
				Type:     "move",
				PlayerId: v.state.PlayerId,
				Action:   action,
			}
			return v, SendMoveCmd(move)
		}
	}
	return v, nil
}

func SendMoveCmd(move message.Move) tea.Cmd {
	return func() tea.Msg {
		return move
	}
}
