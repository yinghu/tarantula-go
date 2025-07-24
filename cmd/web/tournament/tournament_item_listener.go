package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *TournamentService) OnUpdated(kv item.KVUpdate) {
	core.AppLog.Printf("Item update call %v \n", kv)
}
