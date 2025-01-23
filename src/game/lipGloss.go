package game

import "github.com/charmbracelet/lipgloss"

var (
	Red    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Grey   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	Green  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	Black  = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))
	Yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	Blue   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

func colorMessage(success bool, message string) string {
	if success {
		return Green.Render(message)
	}
	return Red.Render(message)
}
