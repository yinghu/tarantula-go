package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type CSQueryer struct {
	*AdminService
}

func (s *CSQueryer) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}

func (s *CSQueryer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	topic := r.PathValue("topic")
	mc, existed := bootstrap.TopicFactoryRegistry[topic]
	if !existed {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: fmt.Sprintf("topic %s not existed", topic)}))
		return
	}
	mf := mc()
	me := mf.Query()
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	defer close(me.QCc())
	req := core.DataRequest{Opt: core.QUERY_DATA_REQUEST, Criteria: me}
	req.Async = me.QCc()
	go s.Cluster().Request(req)
	ms := make([]*protocol.Topic, 0)
	for c := range me.QCc() {
		if !c.Remaining {
			break
		}
		resp, ok := c.Data.(*protocol.Response)

		if ok && resp.Successful {
			for _, data := range resp.Data.List {
				me, err := mf.Topic(data.Value)
				if err != nil {
					continue
				}
				core.AppLog.Debug().Msgf("topic : %v", me)
				ms = append(ms, me)
			}
		}
	}
	w.Write(util.ToJson(ms))
}
