package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *PresenceService) OnRegister(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
	a.LoginReward = conf
}
func (a *PresenceService) OnRelease(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
}
