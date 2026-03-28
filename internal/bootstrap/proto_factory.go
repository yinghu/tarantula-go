package bootstrap

import (
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
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
