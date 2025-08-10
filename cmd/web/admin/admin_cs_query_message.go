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
func (s *CSQueryer) query(query event.Query) {
	s.AdminService.GetJsonAsync(fmt.Sprintf("%s%s/%d", "http://postoffice:8080/postoffice/query/", query.Tag, query.Limit), query.Cc)
}
func (s *CSQueryer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var me event.Query
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	listener := make(chan event.Chunk, 3)
	me.Cc = listener
	defer close(listener)
	go s.query(me)
	for c := range listener {
		if len(c.Data) > 0 {
			w.Write(c.Data)
		}
		if !c.Remaining {
			break
		}
	}
}
