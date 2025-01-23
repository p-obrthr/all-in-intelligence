package game

import (
	"fmt"
	"src/services"
)

type Round struct {
	Deck           Deck
	Players        []Player
	Board          []Card
	Pot            int
	LastRaise      int
	CurrentPlayer  int
	RoundStage     int // 0 = Pre-Flop, 1 = Flop, 2 = Turn, 3 = River
	MsgLog         []string
	BigBlindAmount int
	SmallBlindId   int
	BigBlindId     int
	OpenAIClient   services.OpenAIClient
}

func newRound(Players []Player, client services.OpenAIClient) Round {
	newRound := Round{
		Deck:           newDeck(),
		Players:        Players,
		Board:          make([]Card, 0),
		Pot:            0,
		LastRaise:      0,
		CurrentPlayer:  0,
		RoundStage:     0,
		BigBlindAmount: 40,
		SmallBlindId:   0,
		BigBlindId:     1,
		MsgLog:         []string{},
		OpenAIClient:   client,
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

func (game *Game) endRound() {
	round := &game.Rounds[game.CurrentRound]
	winner := round.determineWinner()
	if winner != nil {
		winner.Money += round.Pot
		message := fmt.Sprintf("%s wins the round with a pot of %d.", winner.Name, round.Pot)
		round.MsgLog = append(round.MsgLog, message)
		for _, player := range round.Players {
			if !player.IsOut {
				round.MsgLog = append(round.MsgLog, fmt.Sprintf("%s cards: %s", player.Name, getCards(player.Cards)))
			}
		}
		round.Pot = 0
	} else {
	}

	game.WaitingForNextRound = true
}

func (round *Round) nextStage() bool {
	switch round.RoundStage {
	case 0:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "Flop.")
	case 1:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "River.")
	case 2:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "Turn.")
	case 3:
		return false
	}
	resetPlayerActions(round)
	round.RoundStage++
	return true
}
