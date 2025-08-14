package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/event"
)

type Tournament interface {
	Join(join event.TournamentEvent) (event.TournamentEvent, error)
	Score(score event.TournamentEvent) (event.TournamentEvent, error)
	Board(update event.TournamentEvent)
}

func main() {
	bootstrap.AppBootstrap(&TournamentService{})
}
