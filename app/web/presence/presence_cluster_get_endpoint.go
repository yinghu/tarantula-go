package main

import (
	"io"
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
	stream, err := s.Cluster().Request(req)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	var data any
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Warn().Msgf("streaming error %s", err.Error())

			break
		}
		if resp.Successful {
			//resp.Data
		}
	}
	w.Write(util.ToJson(data))
}
