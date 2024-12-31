package main

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
}

func newPlayer(i int) Player {

	newPlayer := Player{
		Id:    i,
		Name:  fmt.Sprintf("Player %d", i),
		Cards: make([]Card, 2),
		Money: 1000,
		IsOut: false,
	}

	return newPlayer
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

func (round *Round) raise() (bool, string) {
	amount := 100 + round.LastRaise
	player := &round.Players[round.CurrentPlayer]

	diff := amount - player.InPot

	newBalance := player.Money - diff
	if newBalance >= 0 {
		player.Money = newBalance
		round.Pot += amount
		player.InPot += diff
		round.LastRaise = player.InPot
		return true, fmt.Sprintf("%s raised to %d.", player.Name, player.InPot)
	}
	return false, fmt.Sprintf("%s cannot raise.", player.Name)
}

func (round *Round) fold() (bool, string) {
	currentPlayer := &round.Players[round.CurrentPlayer]
	return true, fmt.Sprintf("%s folded.", currentPlayer.Name)
}

func (round *Round) check() (bool, string) {

	currentPlayer := round.Players[round.CurrentPlayer]

	if currentPlayer.InPot == round.LastRaise {
		return true, fmt.Sprintf("%s checked.", currentPlayer.Name)
	} else {
		return false, fmt.Sprintf("%s cannot check.", currentPlayer.Name)
	}
}

func (round *Round) call() (bool, string) {
	player := &round.Players[round.CurrentPlayer]
	diff := round.LastRaise - player.InPot

	if diff > 0 && player.Money >= diff {
		player.Money -= diff
		round.Pot += diff
		player.InPot += diff
		return true, fmt.Sprintf("%s called to %d.", player.Name, player.InPot)
	}
	return false, fmt.Sprintf("%s cannot call.", player.Name)
}

func checkStatus(player *Player, board *[]Card) {
	cards := append(player.Cards[:], (*board)...)
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
		smallBlindPlayer.InPot += smallBlindAmount
		round.Pot += smallBlindAmount
		round.LastRaise = smallBlindAmount
		round.MsgLog = append(round.MsgLog, fmt.Sprintf("%s paid small blind of %d.", smallBlindPlayer.Name, smallBlindAmount))
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
		bigBlindPlayer.InPot += bigBlindAmount
		round.Pot += bigBlindAmount
		round.LastRaise = bigBlindAmount
		round.MsgLog = append(round.MsgLog, fmt.Sprintf("%s paid big blind of %d.", bigBlindPlayer.Name, bigBlindAmount))
	} else {
	}

	round.CurrentPlayer = (round.BigBlindId + 1) % len(round.Players)
	for round.Players[round.CurrentPlayer].IsOut {
		round.CurrentPlayer = (round.CurrentPlayer + 1) % len(round.Players)
	}
}

func (round *Round) rotateBlinds() {
	round.SmallBlindId = (round.SmallBlindId + 1) % len(round.Players)
	for round.Players[round.SmallBlindId].IsOut {
		round.SmallBlindId = (round.SmallBlindId + 1) % len(round.Players)
	}

	round.BigBlindId = (round.BigBlindId + 1) % len(round.Players)
	for round.Players[round.BigBlindId].IsOut {
		round.BigBlindId = (round.BigBlindId + 1) % len(round.Players)
	}

	round.applyBlinds()
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

func allPlayersHaveActed(round Round) bool {
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
