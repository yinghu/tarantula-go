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

type PostofficeQueryer struct {
	*PostofficeService
}

func (s *PostofficeQueryer) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *PostofficeQueryer) query(query event.Query) {
	buff := persistence.BufferProxy{}
	buff.NewProxy(100)
	query.QCriteria(&buff)
	buff.Flip()
	stat := event.StatEvent{Tag: query.QTag(), Name: event.STAT_TOTAL}
	err := s.Ds.Load(&stat)
	if err != nil {
		query.QCc() <- event.Chunk{Remaining: false, Data: []byte("{\"list\":[]}")}
		return
	}
	query.QCc() <- event.Chunk{Remaining: true, Data: []byte("{\"list\":[")}
	mc := stat.Count
	lmt := query.QLimit()
	s.Ds.List(&buff, func(k, v core.DataBuffer, rev uint64) bool {
		lmt--
		mc--
		cid, _ := v.ReadInt32()
		e := event.CreateEvent(int(cid), nil)
		if e == nil {
			return true
		}
		e.ReadKey(k)
		e.Read(v)
		e.OnRevision(rev)
		ret := util.ToJson(e)
		query.QCc() <- event.Chunk{Remaining: true, Data: ret}
		if lmt > 0 && mc > 0 {
			query.QCc() <- event.Chunk{Remaining: true, Data: []byte(",")}
		}
		return lmt > 0 && mc > 0
	})
	query.QCc() <- event.Chunk{Remaining: false, Data: []byte("]}")}
}

func (s *PostofficeQueryer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	qid, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	me := event.CreateQuery(int32(qid))
	err = json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	defer close(me.QCc())
	go s.query(me)
	for c := range me.QCc() {
		n, err := w.Write(c.Data)
		if err != nil {
			core.AppLog.Printf("Write error %s Num : %d\n", err.Error(), n)
		}
		if !c.Remaining {
			break
		}
	}
}
