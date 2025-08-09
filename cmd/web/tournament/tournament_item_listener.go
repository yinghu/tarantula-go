package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *TournamentService) OnRegister(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d %s\n", conf.Id,conf.Category)
}
func (a *TournamentService) OnRelease(conf item.Configuration) {
	core.AppLog.Printf("item release %d %s\n", conf.Id,conf.Category)
}
