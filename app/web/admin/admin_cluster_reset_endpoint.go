package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminClusterReset struct {
	*AdminService
}

func (s *AdminClusterReset) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *AdminClusterReset) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var co bootstrap.ClusterObject
	err := json.NewDecoder(r.Body).Decode(&co)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	k, v, err := core.Export(&co, 200)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	req := core.DataRequest{Key: k, Value: v, Opt: core.RESET_DATA_REQUEST}
	req.FactoryId = co.FactoryId()
	req.ClassId = co.ClassId()
	req.Mutable = co.Mutable
	aq := make(chan core.Chunk, 3)
	req.Async = aq
	defer close(aq)
	s.Cluster().Request(req)
	for c := range aq {
		if !c.Remaining {
			break
		}
		core.AppLog.Debug().Msgf("payload %v", c.Data)
		w.Write(util.ToJson(c.Data))
	}
}
