package bootstrap

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageEventFactory struct {
	event.ProtoTopicFactoryObj
}

func (p *MessageEventFactory) FromMessageEvent(e event.MessageEvent) (*protocol.Topic, error) {
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

func (p *MessageEventFactory) Message(topic *protocol.Topic) (any, error) {
	me := protocol.MessageEvent{}
	err := anypb.UnmarshalTo(topic.Event.Message, &me, proto.UnmarshalOptions{})
	if err != nil {
		return &me, err
	}
	return &me, nil
}

func (p *MessageEventFactory) Import(criteria []byte) (core.Query, error) {
	q := MessageEventQuery{}
	err := event.Import(&q, criteria, core.COMPOSIT_KEY_MAX)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func (p *MessageEventFactory) Query() core.Query {
	q := MessageEventQuery{}
	q.ClassId = 3
	q.FactoryId = 1
	q.Topic = "message"
	q.Cc = make(chan core.Chunk, 3)
	return &q
}
