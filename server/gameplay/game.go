package gameplay

import (
	"os"
	"server/services"
)

type Game struct {
	Players             []Player
	Rounds              []Round
	CurrentRound        int
	Quit                bool
	BigBlindAmount      int
	StartingMoney       int
	WaitingForNextRound bool
	OpenAIClient        services.OpenAIClient
}

var Id = 1

func NewGame(llmPlayers int, startMoney int, bbAmount int) *Game {

	apiKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists {
		apiKey = ""
	}
	systemPrompt := "You are a professional poker player. You will be given a situation, and you must decide your action. Respond only in JSON format using the following options: { 'action': 'r' } for Raise, { 'action': 'f' } for Fold, { 'action': 'c' } for Call, and { 'action': '_' } for Check."
	client := services.NewOpenAIClient(apiKey, systemPrompt)

	game := &Game{
		BigBlindAmount: bbAmount,
		OpenAIClient:   *client,
		CurrentRound:   -1,
	}

	game.StartingMoney = startMoney

	newMainPlayer := newPlayer(
		1,
		startMoney,
		false,
	)

	game.Players = append(game.Players, newMainPlayer)

	for i := 2; i <= llmPlayers+1; i++ {
		newLlmPlayer := newPlayer(
			i,
			startMoney,
			true,
		)
		game.Players = append(game.Players, newLlmPlayer)
		Id++
	}

	game.StartNewRound()

	// roundJSON, err := json.MarshalIndent(&game.Rounds[game.CurrentRound], "", "  ")

	// if err != nil {
	// 	fmt.Println("Fehler beim Serialisieren des JSON:", err)
	// }

	// fmt.Println("Intended JSON:", string(roundJSON))
	return game
}

func (game *Game) AppendNewPlayer() int {
	Id++
	newMainPlayer := newPlayer(
		Id,
		game.StartingMoney,
		false,
	)
	game.Players = append(game.Players, newMainPlayer)
	return Id
}

func (game *Game) StartNewRound() {
	if game.Quit || game.CurrentRound >= len(game.Rounds) {
		return
	}

	newRound := newRound(game.Players, game.BigBlindAmount)
	game.Rounds = append(game.Rounds, newRound)
	game.CurrentRound++
}
