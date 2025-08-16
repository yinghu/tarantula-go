package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type PostofficeRecoverer struct {
	*PostofficeService
}

func (s *PostofficeRecoverer) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *PostofficeRecoverer) recover(query event.Query) {
	buff := persistence.BufferProxy{}
	buff.NewProxy(100)
	query.QCriteria(&buff)
	buff.Flip()
	stat := event.StatEvent{Tag: query.QTag(), Name: event.STAT_TOTAL}
	err := s.Ds.Load(&stat)
	if err != nil {
		return
	}
	mc := stat.Count
	lmt := 0
	s.Ds.List(&buff, func(k, v core.DataBuffer, rev uint64) bool {
		lmt++
		cid, _ := v.ReadInt32()
		e := event.CreateEvent(int(cid), nil)
		if e == nil {
			return true
		}
		e.ReadKey(k)
		e.Read(v)
		e.OnRevision(rev)
		return true
	})
	core.AppLog.Printf("Total %d recovered from %d\n", lmt, mc)
}

func (s *PostofficeRecoverer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	qid, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	me := event.CreateQuery(int32(qid))
	defer close(me.QCc())
	err = json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	go s.recover(me)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "recovery issued"}))
}
