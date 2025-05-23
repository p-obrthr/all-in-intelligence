package gameplay

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Player struct {
	Id       int
	Name     string
	Money    int
	Cards    []Card
	Status   Status
	InPot    int
	IsOut    bool
	HasActed bool

	IsLLM bool
}

func newPlayer(i int, startingMoney int, isLlm bool) Player {
	name := fmt.Sprintf("Player %d", i)
	if isLlm {
		name += " (LLM)"
	}
	fmt.Println("created " + name)
	return Player{
		Id:    i,
		Name:  name,
		Cards: make([]Card, 2),
		Money: startingMoney,
		InPot: 0,
		IsOut: false,
		IsLLM: isLlm,
	}
}

type Status struct {
	Cards    []Card
	Type     int
	TypeName string
	Score    int
}

func sortPlayersByRanking(players *[]Player) {
	sort.Slice(*players, func(i, j int) bool {
		if (*players)[i].Status.Type != (*players)[j].Status.Type {
			return (*players)[i].Status.Type > (*players)[j].Status.Type
		}
		return (*players)[i].Status.Score > (*players)[j].Status.Score
	})
}

func (round *Round) Raise(amount_optional ...int) (bool, string) {
	var amount int
	if len(amount_optional) > 0 {
		amount = amount_optional[0]
	} else {
		amount = round.LastRaise + 100
	}

	player := round.GetPlayerById(round.CurrentPlayerId)

	diff := amount - player.InPot

	if player.Money >= diff {
		player.Money -= diff
		round.Pot += diff
		player.InPot += diff
		round.LastRaise = amount
		resetPlayerActions(round)
		return true, fmt.Sprintf("%s raised to %d.", player.Name, round.LastRaise)
	} else {
		return false, fmt.Sprintf("%s cannot raise to %d.", player.Name, diff)
	}
}

func (round *Round) Fold() (bool, string) {
	currentPlayer := round.GetPlayerById(round.CurrentPlayerId)
	currentPlayer.IsOut = true
	return true, fmt.Sprintf("%s folded.", currentPlayer.Name)
}

func (round *Round) Check() (bool, string) {
	player := round.GetPlayerById(round.CurrentPlayerId)

	if player.InPot == round.LastRaise {
		return true, fmt.Sprintf("%s checked.", player.Name)
	} else {
		return false, fmt.Sprintf(
			"%s cannot check. Current raise is up to %d.",
			player.Name,
			round.LastRaise,
		)
	}
}

func (round *Round) Call() (bool, string) {
	player := round.GetPlayerById(round.CurrentPlayerId)
	diff := round.LastRaise - player.InPot

	if diff > 0 {
		if player.Money >= diff {
			player.Money -= diff
			round.Pot += diff
			player.InPot += diff
			return true, fmt.Sprintf(
				"%s called to %d.",
				player.Name,
				player.InPot,
			)
		} else {
			round.Pot += player.Money
			player.InPot += player.Money
			player.Money = 0
			return true, fmt.Sprintf(
				"%s is all-in",
				player.Name,
			)
		}
	}
	return true, fmt.Sprintf("%s checked.", player.Name)
}

func (player *Player) checkStatus(board []Card) {
	cards := append(player.Cards, board...)
	cards = filterEmptyCards(cards)

	sortCardsByRanking(&cards)

	ranking, rankedCards := checkRanking(cards)

	score := 0
	for _, card := range rankedCards {
		score += getCardRankValue(card)
	}

	player.Status = Status{
		Cards:    rankedCards,
		Type:     getRankingType(ranking),
		TypeName: ranking,
		Score:    score,
	}

	fmt.Println(player.Id)
	fmt.Println(cards)
	fmt.Println(player.Status)
}

func (round *Round) applyBlinds() {
	smallBlindPlayer := &round.Players[round.SmallBlindId]
	for smallBlindPlayer.IsOut {
		round.SmallBlindId = (round.SmallBlindId + 1) % len(round.Players)
		smallBlindPlayer = &round.Players[round.SmallBlindId]
	}

	smallBlindAmount := round.BigBlindAmount / 2
	if smallBlindPlayer.Money >= smallBlindAmount {
		smallBlindPlayer.Money -= smallBlindAmount
		smallBlindPlayer.InPot = smallBlindAmount
		round.Pot += smallBlindAmount
		round.LastRaise = smallBlindPlayer.InPot
		round.MsgLog = append(
			round.MsgLog,
			fmt.Sprintf(
				"%s paid small blind of %d.",
				smallBlindPlayer.Name,
				smallBlindAmount,
			),
		)
	} else {
	}

	bigBlindPlayer := &round.Players[round.BigBlindId]
	for bigBlindPlayer.IsOut {
		round.BigBlindId = (round.BigBlindId + 1) % len(round.Players)
		bigBlindPlayer = &round.Players[round.BigBlindId]
	}

	bigBlindAmount := round.BigBlindAmount
	if bigBlindPlayer.Money >= bigBlindAmount {
		bigBlindPlayer.Money -= bigBlindAmount
		bigBlindPlayer.InPot = bigBlindAmount
		round.Pot += bigBlindAmount
		round.LastRaise = bigBlindPlayer.InPot
		round.MsgLog = append(
			round.MsgLog,
			fmt.Sprintf(
				"%s paid big blind of %d.",
				bigBlindPlayer.Name,
				bigBlindAmount,
			),
		)
	} else {
	}

	round.CurrentPlayerId = round.findNextPlayerId(bigBlindPlayer.Id)
}

func (round *Round) nextPlayer() {
	totalPlayers := len(round.Players)

	currentPos := -1
	for i, player := range round.Players {
		if player.Id == round.CurrentPlayerId {
			currentPos = i
			break
		}
	}

	if currentPos == -1 {
		fmt.Println("err current player")
		return
	}

	nextPos := (currentPos + 1) % totalPlayers
	for round.Players[nextPos].IsOut {
		nextPos = (nextPos + 1) % totalPlayers
	}

	round.CurrentPlayerId = round.Players[nextPos].Id
}

func (game *Game) CheckPlay() {
	round := &game.Rounds[game.CurrentRound]
	player := round.GetPlayerById(round.CurrentPlayerId)
	if player.IsLLM && !game.WaitingForNextRound {
		isValid := false

		for !isValid {
			action := game.GetLLMAction(*round)
			isValid = game.ProcessPlayerAction(action)

		}
	}
}

type ActionResponse struct {
	Action string `json:"action"`
}

func (game *Game) GetLLMAction(round Round) string {
	prompt := GetPrompt(round)
	response := game.SendRequest(prompt, round)
	return response
}

func GetPrompt(round Round) string {
	player := round.GetPlayerById(round.CurrentPlayerId)
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

func (game *Game) SendRequest(prompt string, round Round) string {
	response, err := game.OpenAIClient.SendMessage(prompt)
	if err != nil {
		return ""
	}
	var actionResponse ActionResponse
	if err := json.Unmarshal([]byte(response), &actionResponse); err != nil {
		return ""
	}
	fmt.Println("-----" + actionResponse.Action)
	return actionResponse.Action
}

func (game *Game) ProcessPlayerAction(message string) bool {
	round := &game.Rounds[game.CurrentRound]
	player := round.GetPlayerById(round.CurrentPlayerId)

	var resp string
	var b bool
	switch message {
	case "r":
		b, resp = round.Raise()
	case "c":
		b, resp = round.Call()
	case "f":
		b, resp = round.Fold()
	case "_":
		b, resp = round.Check()
	}

	if !b {
		return false
	}

	// player.checkStatus(round.Board)

	round.MsgLog = append(round.MsgLog, resp)

	player.HasActed = true
	if getCountActivePlayers(*round) < 2 {
		game.endRound()
		return true
	}

	round.nextPlayer()
	if isAllPlayersHaveActed(*round) && !round.nextStage() {
		game.endRound()
	}
	return true
}

func getCountActivePlayers(round Round) int {
	activePlayers := 0
	for _, p := range round.Players {
		if !p.IsOut {
			activePlayers++
		}
	}
	return activePlayers
}

func isAllPlayersHaveActed(round Round) bool {
	for _, player := range round.Players {
		if !player.IsOut && !player.HasActed {
			return false
		}
	}
	return true
}

func resetPlayerActions(round *Round) {
	for i := range round.Players {
		if !round.Players[i].IsOut {
			round.Players[i].HasActed = false
		}
	}
}
