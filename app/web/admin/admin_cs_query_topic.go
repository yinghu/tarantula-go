package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type CSQueryTopic struct {
	*AdminService
}

func (s *CSQueryTopic) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}

func (s *CSQueryTopic) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	topic := r.PathValue("topic")
	mc, existed := bootstrap.QueryFactoryRegistry[topic]
	if !existed {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: fmt.Sprintf("topic %s not existed", topic)}))
		return
	}
	mf, ok := mc().(core.ProtoTopicFactory)
	if !ok {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: "proto factory not existed"}))
		return
	}
	me := mf.Query()
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	stream, err := s.Cluster().List(me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	ms := make([]any, 0)
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
			for _, data := range resp.Data.List {
				me, err := mf.Topic(data.Value)
				if err != nil {
					continue
				}
				core.AppLog.Debug().Msgf("topic : %v", me)
				e, err := mf.Message(me)
				if err != nil {
					continue
				}
				ms = append(ms, e)
			}
		}
	}
	w.Write(util.ToJson(ms))
}
