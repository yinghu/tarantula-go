package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *TournamentService) OnRegister(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d %s\n", conf.Id, conf.Category)
	if conf.Category == "InstanceSchedule" {
		a.scheduleInstance(conf)
		return
	}
	if conf.Category == "SegmentSchedule" {
		a.scheduleSegment(conf)
		return
	}
	core.AppLog.Printf("Schedule type not supported %s\n", conf.Category)
}
func (a *TournamentService) OnRelease(conf item.Configuration) {
	a.releaseTournament(conf.Id)
}
