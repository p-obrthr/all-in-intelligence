package main

import (
	"sort"
)

type Player struct {
	Name   string
	Money  int
	Cards  [2]Card
	Status Status
}

type Status struct {
	Cards []Card
	Type  int
	Score int
}

func sortPlayersByRanking(players *[]Player) {
	sort.Slice(*players, func(i, j int) bool {
		if (*players)[i].Status.Type != (*players)[j].Status.Type {
			return (*players)[i].Status.Type > (*players)[j].Status.Type
		}
		return (*players)[i].Status.Score > (*players)[j].Status.Score
	})
}
