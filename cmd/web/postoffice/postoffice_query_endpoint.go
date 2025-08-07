package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

type PostofficeQueryer struct {
	*PostofficeService
}

func (s *PostofficeQueryer) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *PostofficeQueryer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tag := r.PathValue("tag")
	limit, err := strconv.ParseInt(r.PathValue("limit"), 10, 32)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	buff := persistence.BufferProxy{}
	buff.NewProxy(100)
	buff.WriteString(tag)
	buff.Flip()
	s.Ds.List(&buff, func(k, v core.DataBuffer) bool {
		limit--
		cid, _ := v.ReadInt32()
		rev, _ := v.ReadInt64()
		core.AppLog.Printf("CID : %d REV : %d", cid, rev)
		return limit == 0
	})
	w.Write(util.ToJson(core.OnSession{Successful: false, Message: tag}))
}
