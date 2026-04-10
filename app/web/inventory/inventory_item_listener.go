package main

import (
	"gameclustering.com/internal/core"
)

func (a *InventoryService) OnRegister(conf core.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)

}
func (a *InventoryService) OnRelease(conf core.Configuration) {
	core.AppLog.Printf("item release %d\n", conf.Id)
}
