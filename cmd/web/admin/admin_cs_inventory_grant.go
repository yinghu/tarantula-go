package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type CSGranter struct {
	*AdminService
}

func (s *CSGranter) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CSGranter) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var me item.OnInventory
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	err = s.ItemService().Manager().Grant(me)
	//id, err := s.Sequence().Id()
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	//me.OnOId(id)
	//me.Source = s.Context()
	//me.DateTime = time.Now()
	//me.OnTopic("message")
	//err = s.Send(&me)
	//if err != nil {
	//w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
	//return
	//}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "granted"}))
}
