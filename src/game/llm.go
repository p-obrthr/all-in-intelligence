package game

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ActionResponse struct {
	Action string `json:"action"`
}

func GetLLMAction(round Round) string {
	prompt := GetPrompt(round)
	response := SendRequest(prompt, round)
	return response
}

func GetPrompt(round Round) string {
	player := &round.Players[round.CurrentPlayer]
	return fmt.Sprintf(
		"Pot: %d\n\nBoard: %s\n\nYour Cards: %s\n\n Your Player Name: %s\n\nYour Money: %d\n\nLog: %s\n\n",
		round.Pot,
		round.Board,
		player.Cards,
		player.Name,
		player.Money,
		strings.Join(round.MsgLog, "\n"),
	)
}

func SendRequest(prompt string, round Round) string {
	response, err := round.OpenAIClient.SendMessage(prompt)
	if err != nil {
		return ""
	}
	var actionResponse ActionResponse
	if err := json.Unmarshal([]byte(response), &actionResponse); err != nil {
		return ""
	}
	return actionResponse.Action
}
