package main

import (
	"gameclustering.com/internal/core"
)

func (a *TournamentService) OnRegister(conf core.Configuration) {
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
func (a *TournamentService) OnRelease(conf core.Configuration) {
	a.releaseTournament(conf.Id)
}
