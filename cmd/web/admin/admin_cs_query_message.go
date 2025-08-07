package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type CSQueryer struct {
	*AdminService
}

func (s *CSQueryer) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CSQueryer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var me event.Query
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	resp := s.AdminService.PostJsonSync(fmt.Sprintf("%s%s/%d", "http://postoffice:8080/postoffice/query/", me.Tag, me.Limit), me)
	w.Write(util.ToJson(resp))
}
