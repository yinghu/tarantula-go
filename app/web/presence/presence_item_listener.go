package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

func (a *PresenceService) OnRegister(conf core.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
	a.LoginReward = conf
	env := util.GitCurBranch()
	err := a.ItemService().SaveRegistration(core.ConfigRegistration{ItemId: conf.Id, App: a.Context(), Env: env.Message})
	if err != nil {
		core.AppLog.Printf("register error %s\n", err.Error())
	}
}
func (a *PresenceService) OnRelease(conf core.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
	env := util.GitCurBranch()
	err := a.ItemService().DeleteRegistration(conf.Id, a.Context(), env.Message)
	if err != nil {
		core.AppLog.Printf("release error %s\n", err.Error())
	}
}
