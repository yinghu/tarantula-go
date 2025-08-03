package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type CSMessager struct {
	*AdminService
}

func (s *CSMessager) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CSMessager) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var me event.MessageEvent
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	resp := s.AdminService.PostJsonSync("http://postoffice:8080/postoffice/publish/"+me.Title, util.ToJson(me))
	w.Write(util.ToJson(resp))
}
