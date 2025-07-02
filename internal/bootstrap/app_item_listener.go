package bootstrap

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

type AppItemListener struct {
	TarantulaService
}

func (a *AppItemListener) OnEnum(e item.Enum) {
	core.AppLog.Printf("%s %v\n", "enum call", e)
}

func (a *AppItemListener) OnCategory(c item.Category) {
	core.AppLog.Printf("%s\n", "category call")
}

func (a *AppItemListener) OnConfiguration(c item.Configuration) {
	core.AppLog.Printf("%s\n", "configuration call")
}
