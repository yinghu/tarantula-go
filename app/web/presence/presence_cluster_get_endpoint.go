package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
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
	//req := core.DataRequest{Key: k, Opt: core.GET_DATA_REQUEST}
	//req.FactoryId = co.FactoryId()
	//req.ClassId = co.ClassId()
	q := protocol.Request{Opt: core.GET_DATA_REQUEST, Data: &protocol.Data{Key: k, Header: &protocol.Header{FactoryId: co.FactoryId(), ClassId: co.ClassId()}}}
	resp, err := s.Cluster().Request(&q)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(resp))
}
