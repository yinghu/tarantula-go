package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *TournamentService) OnRegister(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
}
func (a *TournamentService) OnRelease(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
}
