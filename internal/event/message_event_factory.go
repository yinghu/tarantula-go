package event

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/anypb"
)

const (
	MESSAGE_EVENT_CID  uint32 = 3
	MESSAGE_TOPIC_NAME string = "message"
)

type MessageEventFactory struct {
	ProtoTopicFactoryObj
}

func (p *MessageEventFactory) FromMessageEvent(e *protocol.MessageEvent) (*protocol.Topic, error) {
	tpx := protocol.Topic{Name: MESSAGE_TOPIC_NAME}
	msg := protocol.Event{Header: &protocol.Header{FactoryId: core.EVENT_FACTORY_ID, ClassId: MESSAGE_EVENT_CID}}
	obj, err := anypb.New(e)
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
	err := Import(&q, criteria, core.COMPOSIT_KEY_MAX)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func (p *MessageEventFactory) Query() core.Query {
	q := MessageEventQuery{}
	q.ClassId = MESSAGE_EVENT_CID
	q.FactoryId = core.EVENT_FACTORY_ID
	q.Topic = MESSAGE_TOPIC_NAME
	//q.Cc = make(chan core.Chunk, 3)
	return &q
}
