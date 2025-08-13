package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/event"
)

type Tournament interface {
	Join(join event.TournamentEvent) error
}

func main() {
	bootstrap.AppBootstrap(&TournamentService{})
}
