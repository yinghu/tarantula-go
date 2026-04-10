package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type CSInventoryLoader struct {
	*AdminService
}

func (s *CSInventoryLoader) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}
func (s *CSInventoryLoader) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var me core.OnInventory
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	stock, err := s.Service().ItemService().InventoryManager().Stock(me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(stock))
}
