package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
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
	ivn := item.OnInventory{}
	err := json.NewDecoder(r.Body).Decode(&ivn)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	conf, err := s.ItemService().InventoryManager().Load(ivn.ItemId)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: true, Message: err.Error()}))
		return
	}
	core.AppLog.Printf("Granting item %d\n", conf.Id)
	ivs := make([]item.Inventory, 0)
	s.ItemService().InventoryManager().Validate(conf, func(prop string, c item.Configuration) {
		core.AppLog.Printf("Validating conf %s %s\n", prop, c.Category)
		cat, err := s.ItemService().InventoryManager().LoadCategory(c.Category)
		if err != nil {
			core.AppLog.Printf("Error %s\n", err.Error())
		} else {
			core.AppLog.Printf("Category :%v\n", cat)
			if cat.Scope == item.GRANTABLE_ITEM {
				ivs = append(ivs, item.Inventory{SystemId: ivn.SystemId, TypeId: c.TypeId, Rechargeable: cat.Rechargeable, Amount: c.Amount(cat), UpdateTime: time.Now(), ItemId: c.Id})
			}
		}
	})
	for i := range ivs {
		id, err := s.updateInventory(ivs[i])
		if err != nil {
			w.Write(util.ToJson(core.OnSession{Successful: true, Message: err.Error()}))
			return
		}
		e := event.InventoryEvent{SystemId: ivn.SystemId, InventoryId: id, ItemId: ivn.ItemId, Source: ivn.Source, Description: ivn.Description, GrantTime: time.Now()}
		go s.sendEvent(&e)
	}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "granted"}))
}
func (s *InventoryGranter) sendEvent(e event.Event) {
	oid, _ := s.Sequence().Id()
	e.OnOId(oid)
	e.OnTopic("inventory")
	s.Send(e)
}
