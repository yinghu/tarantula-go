package main

import (
	"bytes"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type CSMessager struct {
	*AdminService
}

func (s *CSMessager) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}
func (s *CSMessager) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	mf := event.MessageEventFactory{}
	var me protocol.MessageEvent
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r.Body)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	err = protojson.Unmarshal(buf.Bytes(), &me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	core.AppLog.Debug().Msgf("event %v", &me)
	tp, err := mf.FromMessage(&me, mf.Header(event.MESSAGE_EVENT_CID))
	tp.Name = event.MESSAGE_TOPIC_NAME
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	tp.Event.Key.Array = core.ToBytes(s.Sequence())
	tp.NodeId = s.NodeId()
	tp.Tag = s.Context()
	resp, err := s.Cluster().Publish(tp)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: resp.Message}))
	go func() {
		mf := persistence.NewCommodityObjectFactory()
		commodity := protocol.Commodity{Name: "gold", Type: "currency", TypeId: "hard currency", Amount: 12, Rechargeable: true}
		kv, err := mf.FromMessage(&commodity, mf.Header(persistence.COMMODITY_OBJECT_ID))
		kv.Key.Array = s.ToBytes(100)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to request %s", err.Error())
			return
		}
		tb := persistence.NewTaskBuilder(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "register"})
		vb := tb.Validator(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "validator"})
		vb.Transaction().Meta(&protocol.Meta{Name: "register"}).Object(kv).Build()
		jb := tb.Job(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "job"})
		jb.Transaction().Meta(&protocol.Meta{Name: "grant"}).Object(kv).Build()
		jb.Build()
		rp, err := s.Cluster().Issue(tb.Build())
		if err != nil {
			core.AppLog.Debug().Msgf("TASK ERR %s", err.Error())
			return
		}
		core.AppLog.Debug().Msgf("TASK %v", rp)

	}()
}
