package main

import (
	"fmt"
	"os"

	"src/game"
	"src/services"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		// log.Fatal("no apikey found in envionment variable")
	}
	systemPrompt := "You are a professional poker player. You will be given a situation, and you must decide your action. Respond only in JSON format using the following options: { 'action': 'r' } for Raise, { 'action': 'f' } for Fold, { 'action': 'c' } for Call, and { 'action': '_' } for Check."
	client := services.NewOpenAIClient(apiKey, systemPrompt)

	game := game.NewGame(*client)
	p := tea.NewProgram(game)
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
