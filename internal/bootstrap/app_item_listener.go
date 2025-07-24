package bootstrap

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

type AppItemListener struct {
	TarantulaService
}

func (a *AppItemListener) OnUpdated(kv item.KVUpdate) {
	core.AppLog.Printf("Item update call %v \n", kv)
}
