package bootstrap

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromMessageEvent(e event.MessageEvent) (*protocol.Topic, error) {
	tpx := protocol.Topic{NodeId: e.NodeId(), Tag: e.Tag(), Name: e.Topic()}
	msg := protocol.Event{Header: &protocol.Header{FactoryId: e.FactoryId(), ClassId: e.ClassId()}, Id: uint64(e.OId())}
	obj, err := anypb.New(&protocol.MessageEvent{Title: e.Title, Message: e.Message, Source: e.Source, DateTime: timestamppb.New(e.DateTime)})
	if err != nil {
		return &tpx, err
	}
	msg.Message = obj
	tpx.Event = &msg
	return &tpx, nil
}

func ToRequest(topic *protocol.Topic) (*protocol.Request, error) {
	req := protocol.Request{Opt: core.CREATE_DATA_REQUEST, Prefix: util.Hash([]byte(topic.Name))}
	if topic.Event.Id <= 0 {
		return &req, fmt.Errorf("id cannot be less than 0")
	}
	buff := core.NewBuffer(8)
	if err := buff.WriteUInt64(topic.Event.Id); err != nil {
		return &req, err
	}
	buff.Flip()
	key, err := buff.Read(0)
	if err != nil {
		return &req, err
	}
	value, err := proto.Marshal(topic)
	if err != nil {
		return &req, err
	}
	data := protocol.Data{Header: topic.Event.Header, Key: key, Value: value}
	req.Data = &data
	return &req, nil
}
