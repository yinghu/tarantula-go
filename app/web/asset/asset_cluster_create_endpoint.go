package main

import (
	"encoding/json"
	"io"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AssetClusterCreate struct {
	*AssetService
}

func (s *AssetClusterCreate) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *AssetClusterCreate) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
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
	req := core.DataRequest{Key: k, Value: v, Opt: core.CREATE_DATA_REQUEST}
	req.FactoryId = co.FactoryId()
	req.ClassId = co.ClassId()
	req.Mutable = co.Mutable

	stream, err := s.Cluster().Request(req)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	ret := core.OnSession{}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Warn().Msgf("streaming error %s", err.Error())
			ret.Successful = false
			ret.Message = err.Error()
			break
		}
		ret.Successful = resp.Successful
		ret.Message = resp.Message
	}
	w.Write(util.ToJson(ret))
}
