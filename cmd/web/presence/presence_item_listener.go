package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *PresenceService) OnUpdated(kv item.KVUpdate) {
	core.AppLog.Printf("Item update call %v \n", kv)
}
