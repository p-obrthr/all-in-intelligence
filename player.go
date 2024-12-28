package main

import (
	"fmt"
	"sort"
)

type Player struct {
	Name   string
	Money  int
	Cards  [2]Card
	Status Status
	InPot  int
	IsOut  bool
}

func newPlayer(i int) Player {
	return Player{
		Name:  fmt.Sprintf("Player %d", i),
		Money: 1000,
		IsOut: false,
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

func (game *Game) doRaise(amount int) bool {
	newBalance := game.Players[game.CurrentPlayer].Money - amount
	if newBalance >= 0 {
		game.Players[game.CurrentPlayer].Money = newBalance
		game.Pot += amount
		game.Players[game.CurrentPlayer].InPot += amount
		game.LastRaise = amount
		return true
	} else {
		return false
	}
}

func (game *Game) doFold() bool {
	if !game.Players[game.CurrentPlayer].IsOut {
		game.Players[game.CurrentPlayer].IsOut = true
		return true
	} else {
		return false
	}
}

func (game *Game) doCall() bool {
	diff := game.LastRaise - game.Players[game.CurrentPlayer].InPot

	if diff > 0 {
		if game.Players[game.CurrentPlayer].Money >= diff {
			game.Players[game.CurrentPlayer].Money -= diff
			game.Pot += diff
			game.Players[game.CurrentPlayer].InPot += diff
			return true
		} else {
			return false
		}
	}

	return true
}

func checkStatus(player *Player, board *Board) string {
	var cards []Card

	for i := 0; i < len(player.Cards); i++ {
		cards = append(cards, player.Cards[i])
	}

	for i := 0; i < len(board); i++ {
		cards = append(cards, board[i])
	}

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

	return ""
}
