package main

import (
	"gameclustering.com/internal/core"
)

func (a *AssetService) OnRegister(conf core.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
}
func (a *AssetService) OnRelease(conf core.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
}
