package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *MahjongService) OnRegister(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
	
}
func (a *MahjongService) OnRelease(conf item.Configuration) {
	core.AppLog.Printf("item release %d\n", conf.Id)
}
