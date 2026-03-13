package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type PresenceClusterGet struct {
	*PresenceService
}

func (s *PresenceClusterGet) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *PresenceClusterGet) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	kn := r.PathValue("key")

	var co bootstrap.ClusterObject
	co.Key = kn
	k, _, err := core.Export(&co, 200)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	req := core.DataRequest{Key: k, Opt: core.GET_DATA_REQUEST}
	req.FactoryId = co.FactoryId()
	req.ClassId = co.ClassId()
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
	w.Write(util.ToJson(co))
}
