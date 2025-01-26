package main

import (
	"log"
	"os"
	"src/game"
	"src/services"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		apiKey = ""
	}
	systemPrompt := "You are a professional poker player. You will be given a situation, and you must decide your action. Respond only in JSON format using the following options: { 'action': 'r' } for Raise, { 'action': 'f' } for Fold, { 'action': 'c' } for Call, and { 'action': '_' } for Check."
	client := services.NewOpenAIClient(apiKey, systemPrompt)

	game := game.New(*client)
	p := tea.NewProgram(&game, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
