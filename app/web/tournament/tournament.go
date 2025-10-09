package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/event"
)

type Tournament interface {
	Join(join event.TournamentEvent) (event.TournamentEvent, error)
	Score(score event.TournamentEvent) (event.TournamentEvent, error)
	OnBoard(update event.TournamentEvent)
	Listing(query event.TournamentEvent) []RaceEntry
	Start() error
}

func main() {
	bootstrap.AppBootstrap(&TournamentService{})
}
