package event

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/known/anypb"
)

const (
	LOG_EVENT_CID  uint32 = 1
	LOG_TOPIC_NAME string = "log"
)

type LogEventFactory struct {
	ProtoTopicFactoryObj
}

func (p *LogEventFactory) FromLogEvent(e *protocol.LogEvent) (*protocol.Topic, error) {
	tpx := protocol.Topic{Name: LOG_TOPIC_NAME}
	msg := protocol.Event{Header: &protocol.Header{FactoryId: core.EVENT_FACTORY_ID, ClassId: LOG_EVENT_CID}}
	obj, err := anypb.New(e)
	if err != nil {
		return &tpx, err
	}
	msg.Message = obj
	tpx.Event = &msg
	return &tpx, nil
}

func (p *LogEventFactory) Message(topic *protocol.Topic) (any, error) {
	me := protocol.LogEvent{}
	err := anypb.UnmarshalTo(topic.Event.Message, &me, proto.UnmarshalOptions{})
	if err != nil {
		return &me, err
	}
	return &me, nil
}

func (p *LogEventFactory) Import(criteria []byte) (core.Query, error) {
	q := LogEventQuery{}
	err := Import(&q, criteria, core.COMPOSIT_KEY_MAX)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func (p *LogEventFactory) Query() core.Query {
	q := LogEventQuery{}
	q.ClassId = LOG_EVENT_CID
	q.FactoryId = core.EVENT_FACTORY_ID
	q.Topic = LOG_TOPIC_NAME
	q.Cc = make(chan core.Chunk, 3)
	return &q
}
