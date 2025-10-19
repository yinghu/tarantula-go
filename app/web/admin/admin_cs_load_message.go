package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type CSMessageLoader struct {
	*AdminService
}

func (s *CSMessageLoader) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CSMessageLoader) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var me event.MessageEvent
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	q := event.QWithTag{Ee: &me}
	q.Cc = make(chan core.Chunk, 3)
	s.Load(&q)
	defer close(q.QCc())
	for c := range q.QCc() {
		if len(c.Data) > 0 {
			w.Write(c.Data)
		}
		if !c.Remaining {
			break
		}
	}
}
