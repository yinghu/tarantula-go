package bootstrap

import (
	"gameclustering.com/internal/core"
)

type AppItemListener struct {
	TarantulaService
}

func (a *AppItemListener) OnUpdated(kv core.KVUpdate) {
	core.AppLog.Printf("Item update call %v \n", kv)
}
