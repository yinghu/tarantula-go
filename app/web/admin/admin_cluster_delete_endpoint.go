package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminClusterDelete struct {
	*AdminService
}

func (s *AdminClusterDelete) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *AdminClusterDelete) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	kn := r.PathValue("key")

	var co bootstrap.ClusterObject
	co.Key = kn
	k, _, err := core.Export(&co, 200)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	req := core.DataRequest{Key: k, Opt: core.DELETE_DATA_REQUEST}
	req.FactoryId = co.FactoryId()
	req.ClassId = co.ClassId()
	resp, err := s.Cluster().Request(req)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(resp))

}
