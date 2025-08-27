package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type InventoryLoader struct {
	*InventoryService
}

func (s *InventoryLoader) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *InventoryLoader) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ivn := item.OnInventory{}
	err := json.NewDecoder(r.Body).Decode(&ivn)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	list, err := s.loadInventory(ivn)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: true, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(list))
}
