package main

import (
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
	buff.WriteString(query.Tag)
	buff.Flip()
	query.Cc <- event.Chunk{Remaining: true, Data: []byte("[")}
	s.Ds.List(&buff, func(k, v core.DataBuffer) bool {
		query.Limit--
		cid, _ := v.ReadInt32()
		rev, _ := v.ReadInt64()
		e := event.CreateEvent(int(cid), nil)
		if e == nil {
			return true
		}
		e.ReadKey(k)
		e.Read(v)
		core.AppLog.Printf("CID : %d REV : %d %v\n", cid, rev, e)
		query.Cc <- event.Chunk{Remaining: true, Data: util.ToJson(e)}
		if query.Limit > 0 {
			query.Cc <- event.Chunk{Remaining: true, Data: []byte(",")}
		}
		return query.Limit > 0
	})
	query.Cc <- event.Chunk{Remaining: false, Data: []byte("]")}
}

func (s *PostofficeQueryer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	listener := make(chan event.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	tag := r.PathValue("tag")
	limit, err := strconv.ParseInt(r.PathValue("limit"), 10, 32)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	q := event.Query{Tag: tag, Limit: int32(limit)}
	q.Cc = listener
	go s.query(q)
	for c := range listener {
		w.Write(c.Data)
		if !c.Remaining {
			break
		}
	}
}
