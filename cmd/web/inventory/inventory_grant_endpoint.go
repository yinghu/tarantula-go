package main

import (
	"net/http"
	"strconv"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type InventoryGranter struct {
	*InventoryService
}

func (s *InventoryGranter) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *InventoryGranter) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	qid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	err = s.updateInventory(Inventory{SystemId: rs.SystemId, TypeId: "Coin", Rechargeable: true, Amount: 100, UpdateTime: time.Now()}, qid)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: true, Message: err.Error()}))
		return
	}
	e := event.InventoryEvent{SystemId: rs.SystemId, InventoryId: 10, ItemId: qid, Source: "web", Description: "event test", GrantTime: time.Now()}
	oid, _ := s.Sequence().Id()
	e.OnOId(oid)
	e.OnTopic("inventory")
	s.Send(&e)
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "granted"}))
}
