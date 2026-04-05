package event

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/anypb"
)

const (
	LOGIN_EVENT_CID  uint32 = 4
	LOGIN_TOPIC_NAME string = "login"
)

type LoginEventFactory struct {
	ProtoTopicFactoryObj
}

func (p *LoginEventFactory) FromLoginEvent(e *protocol.LoginEvent) (*protocol.Topic, error) {
	tpx := protocol.Topic{Name: LOGIN_TOPIC_NAME}
	msg := protocol.Event{Header: &protocol.Header{FactoryId: core.EVENT_FACTORY_ID, ClassId: LOGIN_EVENT_CID}}
	obj, err := anypb.New(e)
	if err != nil {
		return &tpx, err
	}
	msg.Message = obj
	tpx.Event = &msg
	return &tpx, nil
}

func (p *LoginEventFactory) Message(topic *protocol.Topic) (any, error) {
	me := protocol.LoginEvent{}
	err := anypb.UnmarshalTo(topic.Event.Message, &me, proto.UnmarshalOptions{})
	if err != nil {
		return &me, err
	}
	return &me, nil
}

func (p *LoginEventFactory) Import(criteria []byte) (core.Query, error) {
	q := LoginEventQuery{}
	err := Import(&q, criteria, core.COMPOSIT_KEY_MAX)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func (p *LoginEventFactory) Query() core.Query {
	q := LoginEventQuery{}
	q.ClassId = LOGIN_EVENT_CID
	q.FactoryId = core.EVENT_FACTORY_ID
	q.Topic = LOGIN_TOPIC_NAME
	q.Cc = make(chan core.Chunk, 3)
	return &q
}
