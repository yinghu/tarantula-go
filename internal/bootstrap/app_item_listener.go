package bootstrap

import "gameclustering.com/internal/core"

type AppItemListener struct {
	TarantulaService
}

func (a *AppItemListener) OnUpdate() {
	core.AppLog.Printf("Item update call \n")
}
