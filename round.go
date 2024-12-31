package main

import "fmt"

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
}

func newRound(Players []Player) Round {
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
	}

	for _, player := range Players {
		player.Cards[0] = drawCard(&newRound.Deck)
		player.Cards[1] = drawCard(&newRound.Deck)
	}

	newRound.applyBlinds()

	return newRound
}

func (game *Game) endRound() {
	round := game.Rounds[game.CurrentRound]

	winnerIndex := round.determineWinner()
	if winnerIndex >= 0 {
		winner := &round.Players[winnerIndex]
		round.MsgLog = append(round.MsgLog, fmt.Sprintf("%s wins", winner.Name))
	}

	if len(round.Players) < 2 {
		game.Players[0].Money += round.Pot
		game.Quit = true
	} else {
		for _, player := range game.Players {
			if player.Money == 0 {
				player.IsOut = true
			} else {
				player.IsOut = false
			}
		}
	}

	resetPlayerActions(&game.Rounds[game.CurrentRound])

	game.CurrentRound++
	newRound := newRound(game.Players)
	game.Rounds = append(game.Rounds, newRound)

	newRound.applyBlinds()
}

func (round *Round) nextStage() bool {
	switch round.RoundStage {
	case 0:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "Flop.")
		break
	case 1:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "River.")
		break
	case 2:
		round.Board = append(round.Board, drawCard(&round.Deck))
		round.MsgLog = append(round.MsgLog, "Turn.")
		break
	case 3:
		return false
	}
	resetPlayerActions(round)
	round.RoundStage++
	return true
}
