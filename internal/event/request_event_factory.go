package event

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/anypb"
)

const (
	REQUEST_EVENT_CID  uint32 = 5
	REQUEST_TOPIC_NAME string = "request"
)

type RequestEventFactory struct {
	ProtoTopicFactoryObj
}

func (p *RequestEventFactory) FromRequestEvent(e *protocol.RequestEvent) (*protocol.Topic, error) {
	tpx := protocol.Topic{Name: REQUEST_TOPIC_NAME}
	msg := protocol.Event{Header: &protocol.Header{FactoryId: core.EVENT_FACTORY_ID, ClassId: REQUEST_EVENT_CID}}
	obj, err := anypb.New(e)
	if err != nil {
		return &tpx, err
	}
	msg.Message = obj
	tpx.Event = &msg
	return &tpx, nil
}

func (p *RequestEventFactory) Message(topic *protocol.Topic) (any, error) {
	me := protocol.RequestEvent{}
	err := anypb.UnmarshalTo(topic.Event.Message, &me, proto.UnmarshalOptions{})
	if err != nil {
		return &me, err
	}
	return &me, nil
}

func (p *RequestEventFactory) Import(criteria []byte) (core.Query, error) {
	q := RequestEventQuery{}
	err := Import(&q, criteria, core.COMPOSIT_KEY_MAX)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func (p *RequestEventFactory) Query() core.Query {
	q := RequestEventQuery{}
	q.ClassId = REQUEST_EVENT_CID
	q.FactoryId = core.EVENT_FACTORY_ID
	q.Topic = REQUEST_TOPIC_NAME
	//q.Cc = make(chan core.Chunk, 3)
	return &q
}
