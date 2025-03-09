package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type ConfigView struct {
	focusIndex int
	inputs     []textinput.Model
}

func NewConfigView() *ConfigView {
	v := ConfigView{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range v.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Amount of LLM Players"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		// case 1:
		// 	t.Placeholder = "Amount of real websocket Players to join "
		// 	t.CharLimit = 64
		case 1:
			t.Placeholder = "Small blind amount"
		}

		v.inputs[i] = t
	}

	return &v
}

func (v ConfigView) Init() tea.Cmd {
	return textinput.Blink
}

func (v ConfigView) HandleKey(msg interface{}) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && v.focusIndex == len(v.inputs) {
				players, err := strconv.Atoi(v.inputs[0].Value())
				if err != nil {
					fmt.Println("err converting player amount:", err)
				}

				bbAmount, err := strconv.Atoi(v.inputs[1].Value())
				if err != nil {
					fmt.Println("err converting big blind amount:", err)
				}

				return NewGameView(players, bbAmount), tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				v.focusIndex--
			} else {
				v.focusIndex++
			}

			if v.focusIndex > len(v.inputs) {
				v.focusIndex = 0
			} else if v.focusIndex < 0 {
				v.focusIndex = len(v.inputs)
			}

			cmds := make([]tea.Cmd, len(v.inputs))
			for i := 0; i <= len(v.inputs)-1; i++ {
				if i == v.focusIndex {
					cmds[i] = v.inputs[i].Focus()
					v.inputs[i].PromptStyle = focusedStyle
					v.inputs[i].TextStyle = focusedStyle
					continue
				}
				v.inputs[i].Blur()
				v.inputs[i].PromptStyle = noStyle
				v.inputs[i].TextStyle = noStyle
			}

			return v, tea.Batch(cmds...)

		case "esc":
			return NewStartView(), nil
		}
	}

	cmd := v.updateInputs(msg)

	return v, cmd
}

func (v ConfigView) Render() string {
	return v.GetConfigView()
}

func (v *ConfigView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(v.inputs))

	for i := range v.inputs {
		v.inputs[i], cmds[i] = v.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (v ConfigView) GetConfigView() string {
	var b strings.Builder

	b.WriteString("Lets do some config for the game.\n\n")

	for i := range v.inputs {
		b.WriteString(v.inputs[i].View())
		if i < len(v.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if v.focusIndex == len(v.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
