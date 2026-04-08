package event

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/anypb"
)

const (
	REGISTER_EVENT_CID  uint32 = 2
	REGISTER_TOPIC_NAME string = "register"
)

type RegisterEventFactory struct {
	ProtoTopicFactoryObj
}

func (p *RegisterEventFactory) FromRegisterEvent(e *protocol.RegisterEvent) (*protocol.Topic, error) {
	tpx := protocol.Topic{Name: REGISTER_TOPIC_NAME}
	msg := protocol.Event{Header: &protocol.Header{FactoryId: core.EVENT_FACTORY_ID, ClassId: REGISTER_EVENT_CID}}
	obj, err := anypb.New(e)
	if err != nil {
		return &tpx, err
	}
	msg.Message = obj
	tpx.Event = &msg
	return &tpx, nil
}

func (p *RegisterEventFactory) Message(topic *protocol.Topic) (any, error) {
	me := protocol.RegisterEvent{}
	err := anypb.UnmarshalTo(topic.Event.Message, &me, proto.UnmarshalOptions{})
	if err != nil {
		return &me, err
	}
	return &me, nil
}

func (p *RegisterEventFactory) Import(criteria []byte) (core.Query, error) {
	q := RegisterEventQuery{}
	err := Import(&q, criteria, core.COMPOSIT_KEY_MAX)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func (p *RegisterEventFactory) Query() core.Query {
	q := RegisterEventQuery{}
	q.ClassId = REGISTER_EVENT_CID
	q.FactoryId = core.EVENT_FACTORY_ID
	q.Topic = REGISTER_TOPIC_NAME
	//q.Cc = make(chan core.Chunk, 3)
	return &q
}
