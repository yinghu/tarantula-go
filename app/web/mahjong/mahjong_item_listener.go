package main

import (
	"gameclustering.com/internal/core"
)

func (a *MahjongService) OnRegister(conf core.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)

}
func (a *MahjongService) OnRelease(conf core.Configuration) {
	core.AppLog.Printf("item release %d\n", conf.Id)
}
