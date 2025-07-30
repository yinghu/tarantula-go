package main

import (
	"encoding/json"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *TournamentService) OnUpdated(kv item.KVUpdate) {
	core.AppLog.Printf("Item update call %v \n", kv)
	var reg item.ConfigRegistration
	err := json.Unmarshal([]byte(kv.Value), &reg)
	if err != nil {
		return
	}
	ins, err := a.ItemService().Loader().Load(reg.ItemId)
	if err != nil {
		return
	}
	core.AppLog.Printf("%v\n", ins)
}
