package main

import (
	"bytes"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/types/known/anypb"
)

type CSMessager struct {
	*AdminService
}

func (s *CSMessager) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}
func (s *CSMessager) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
	id, err := s.Sequence().Id()
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	//obj, err := anypb.New(&protocol.MessageEvent{Title: e.Title, Message: e.Message, Source: e.Source, DateTime: timestamppb.New(e.DateTime)})
	obj, err := anypb.New(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	e := protocol.Event{Id: uint64(id), Header: &protocol.Header{FactoryId: 1, ClassId: 3}, Message: obj}
	t := protocol.Topic{Event: &e, NodeId: s.NodeId(), Tag: s.Context(), Name: "message"}
	//me.OnOId(id)
	//me.Source = s.Context()
	//me.DateTime = time.Now()
	//tf := bootstrap.MessageEventFactory{}
	//tp, err := tf.FromMessageEvent(me)
	//if err != nil {
	//w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
	//return
	//}
	err = s.Cluster().Publish(&t)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "message delivered"}))
}
