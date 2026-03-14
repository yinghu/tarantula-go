package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type AdminClusterPull struct {
	*AdminService
}

func (s *AdminClusterPull) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *AdminClusterPull) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ringToken, err := strconv.ParseUint(r.PathValue("ringToken"), 10, 32)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	limit, err := strconv.ParseInt(r.PathValue("limit"), 10, 32)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	query := event.QWithTag{Id: event.RING_PULL_QID, Limit: int32(limit)}
	req := core.DataRequest{Prefix: uint32(ringToken), Opt: core.PULL_DATA_REQUEST, Criteria: &query}
	aq := make(chan core.Chunk, 3)
	req.Async = aq
	defer close(aq)
	s.Cluster().Request(req)
	for c := range aq {
		if !c.Remaining {
			break
		}
		core.AppLog.Debug().Msgf("payload %v", c.Data)
	}
	w.Write(util.ToJson(core.OnSession{Successful: false, Message: "pending"}))
}
