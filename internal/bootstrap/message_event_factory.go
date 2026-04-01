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

type MessageEventFactory struct {
	core.ProtoTopicFactoryObj
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

func (p *MessageEventFactory) Request(topic *protocol.Topic) (*protocol.Request, error) {
	p.Target = topic
	req := protocol.Request{Opt: core.CREATE_DATA_REQUEST}
	if topic.Event.Id <= 0 {
		return &req, fmt.Errorf("id cannot be less than 0")
	}
	buff := core.NewBuffer(core.COMPOSIT_KEY_MAX)
	p.WriteKey(buff)
	buff.Flip()
	key, err := buff.Read(0)
	if err != nil {
		return &req, err
	}
	req.Prefix = util.Hash(key)
	value, err := proto.Marshal(topic)
	if err != nil {
		return &req, err
	}
	data := protocol.Data{Header: topic.Event.Header, Key: key, Value: value}
	req.Data = &data
	return &req, nil
}

func (p *MessageEventFactory) WriteKey(key core.DataBuffer) error {
	me := protocol.MessageEvent{}
	err := anypb.UnmarshalTo(p.Target.Event.Message, &me, proto.UnmarshalOptions{})
	if err != nil {
		return err
	}
	return key.WriteUInt64(p.Target.Event.Id)
}

func (p *MessageEventFactory) Query(criteria []byte) (core.Query, error) {
	q := MessageEventQuery{}
	err := event.Import(&q, criteria, core.COMPOSIT_KEY_MAX)
	if err != nil {
		return &q, err
	}
	return &q, nil
}
