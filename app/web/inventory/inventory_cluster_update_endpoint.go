package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type InventoryClusterUpdate struct {
	*InventoryService
}

func (s *InventoryClusterUpdate) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *InventoryClusterUpdate) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
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
	//req := core.DataRequest{Key: k, Value: v, Opt: core.UPDATE_DATA_REQUEST}
	//req.FactoryId = co.FactoryId()
	//req.ClassId = co.ClassId()
	//req.Revision = co.Rev
	q := protocol.Request{Opt: core.UPDATE_DATA_REQUEST, Data: &protocol.Data{Key: k, Value: v, Header: &protocol.Header{FactoryId: co.FactoryId(), ClassId: co.ClassId(), Revision: co.Rev}}}
	resp, err := s.Cluster().Request(&q)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(resp))
}
