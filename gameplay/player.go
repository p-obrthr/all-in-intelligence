package gameplay

import (
	"fmt"
	"sort"
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
	IsLLM    bool
}

func newPlayer(i int, startingMoney int, isLlm bool) Player {
	name := fmt.Sprintf("Player %d", i)
	if isLlm {
		name += " (LLM)"
	}
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

	player := &round.Players[round.CurrentPlayer]

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
	currentPlayer := &round.Players[round.CurrentPlayer]
	currentPlayer.IsOut = true
	return true, fmt.Sprintf("%s folded.", currentPlayer.Name)
}

func (round *Round) Check() (bool, string) {
	player := &round.Players[round.CurrentPlayer]

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
	player := &round.Players[round.CurrentPlayer]
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

		}
	}
	return true, fmt.Sprintf("%s checked.", player.Name)
}

func checkStatus(player *Player, board *[]Card) {
	cards := append(player.Cards, *board...)
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

	round.CurrentPlayer = (round.BigBlindId + 1) % len(round.Players)
	for round.Players[round.CurrentPlayer].IsOut {
		round.CurrentPlayer = (round.CurrentPlayer + 1) % len(round.Players)
	}
}

func (round *Round) nextPlayer() {
	totalPlayers := len(round.Players)

	for i := 1; i < totalPlayers; i++ {
		nextPlayerId := (round.CurrentPlayer + i) % totalPlayers
		if !round.Players[nextPlayerId].IsOut {
			round.CurrentPlayer = nextPlayerId
			break
		}
	}
}

func (game *Game) ProcessPlayerAction(actionSuccessful bool, message string) {
	round := &game.Rounds[game.CurrentRound]
	player := &round.Players[round.CurrentPlayer]

	round.MsgLog = append(round.MsgLog, message)

	if actionSuccessful {
		player.HasActed = true
		if getCountActivePlayers(*round) < 2 {
			game.endRound()
			return
		}

		round.nextPlayer()
		if isAllPlayersHaveActed(*round) && !round.nextStage() {
			game.endRound()
		}
	}
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
