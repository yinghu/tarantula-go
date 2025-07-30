package main

import (
	"encoding/json"
	"strconv"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (a *TournamentService) OnUpdated(kv item.KVUpdate) {
	itemId, err := strconv.ParseInt(kv.Key, 10, 64)
	if err != nil {
		core.AppLog.Printf("Key should be int64 %s\n", kv.Key)
	}
	if kv.Opt.IsCreate || kv.IsModify {
		var reg item.ConfigRegistration
		err := json.Unmarshal([]byte(kv.Value), &reg)
		if err != nil {
			core.AppLog.Printf("Value should be json format %v\n", kv.Value)
			return
		}
		if reg.ItemId != itemId {
			core.AppLog.Printf("Key not matched %d : %d\n", itemId, reg.ItemId)
			return
		}
		ins, err := a.ItemService().Loader().Load(reg.ItemId)
		if err != nil {
			return
		}
		core.AppLog.Printf("Item registered %d\n", ins.Id)
		return
	}
	core.AppLog.Printf("Item released %d\n", itemId)

}
