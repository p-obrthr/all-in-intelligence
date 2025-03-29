package gameplay

import (
	"fmt"
)

type Round struct {
	Deck            Deck
	Players         []Player
	Board           []Card
	Pot             int
	LastRaise       int
	CurrentPlayerId int
	RoundStage      RoundStage
	MsgLog          []string
	BigBlindAmount  int
	SmallBlindId    int
	BigBlindId      int
}

type RoundStage int

const (
	Preflop RoundStage = iota
	Flop
	Turn
	River
)

func newRound(Players []Player, bigBlindAmount int) Round {
	newRound := Round{
		Deck:            newDeck(),
		Players:         Players,
		Board:           make([]Card, 0),
		Pot:             0,
		LastRaise:       0,
		CurrentPlayerId: 0,
		RoundStage:      Preflop,
		BigBlindAmount:  bigBlindAmount,
		SmallBlindId:    0,
		BigBlindId:      1,
		MsgLog:          []string{},
	}

	for i := range newRound.Players {
		newRound.Players[i].Cards[0] = drawCard(&newRound.Deck)
		newRound.Players[i].Cards[1] = drawCard(&newRound.Deck)
		newRound.Players[i].InPot = 0
		if newRound.Players[i].Money > 0 {
			newRound.Players[i].IsOut = false
		}
	}

	newRound.applyBlinds()

	return newRound
}

func (round *Round) GetPlayerById(id int) *Player {
	for i := range round.Players {
		if round.Players[i].Id == id {
			return &round.Players[i]
		}
	}
	return nil
}

func (round *Round) findNextPlayerId(startId int) int {
	totalPlayers := len(round.Players)

	for i := 0; i < totalPlayers; i++ {
		nextId := (startId + i) % totalPlayers
		if !round.Players[nextId].IsOut {
			return round.Players[nextId].Id
		}
	}

	return round.Players[startId].Id
}

func (game *Game) endRound() {
	round := &game.Rounds[game.CurrentRound]
	winner := round.determineWinner()
	if winner != nil {
		winner.Money += round.Pot
		message := fmt.Sprintf(
			"%s wins the round with a pot of %d.",
			winner.Name,
			round.Pot,
		)
		round.MsgLog = append(round.MsgLog, message)
		for _, player := range round.Players {
			if !player.IsOut {
				round.MsgLog = append(
					round.MsgLog,
					fmt.Sprintf(
						"%s cards: %s",
						player.Name,
						GetCards(player.Cards),
					),
				)
			}
		}
		round.Pot = 0
	} else {
	}

	game.WaitingForNextRound = true
}

func (round *Round) nextStage() bool {
	switch round.RoundStage {
	case Preflop:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "Flop.")
	case Flop:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "Turn.")
	case Turn:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "River.")
	case River:
		return false
	}
	resetPlayerActions(round)
	round.RoundStage++
	for i := 0; i < len(round.Players); i++ {
		round.Players[i].checkStatus(round.Board)
	}
	return true
}

func (round *Round) determineWinner() *Player {
	for i := range round.Players {
		if !round.Players[i].IsOut {
			round.Players[i].checkStatus(round.Board)
		}
	}

	sortPlayersByRanking(&round.Players)

	for i := range round.Players {
		if !round.Players[i].IsOut {
			return &round.Players[i]
		}
	}
	return nil
}
