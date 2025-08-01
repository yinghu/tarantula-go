package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *InventoryService) OnRegister(conf item.Configuration) {
	core.AppLog.Printf("item reigster %d\n", conf.Id)
}
func (a *InventoryService) OnRelease(conf item.Configuration) {
	core.AppLog.Printf("item release %d\n", conf.Id)
}
