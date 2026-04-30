package event

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	LOG_EVENT_CID  uint32 = 1
	LOG_TOPIC_NAME string = "log"

	REGISTER_EVENT_CID  uint32 = 2
	REGISTER_TOPIC_NAME string = "register"

	MESSAGE_EVENT_CID  uint32 = 3
	MESSAGE_TOPIC_NAME string = "message"

	LOGIN_EVENT_CID  uint32 = 4
	LOGIN_TOPIC_NAME string = "login"

	REQUEST_EVENT_CID  uint32 = 5
	REQUEST_TOPIC_NAME string = "request"

	TRANSACTION_EVENT_CID  uint32 = 6
	TRANSACTION_TOPIC_NAME string = "transaction"
)

type MessageTopic func() proto.Message

type ProtoTopicFactoryObj struct {
	Target *protocol.Topic
	core.QueryFactoryObj
	Mt MessageTopic
}

func (p *ProtoTopicFactoryObj) Request(topic *protocol.Topic) (*protocol.Request, error) {
	p.Target = topic
	req := protocol.Request{Opt: core.CREATE_DATA_REQUEST}
	if len(topic.Event.Key.Array) == 0 {
		return &req, fmt.Errorf("key must be assigned")
	}
	value, err := proto.Marshal(topic)
	if err != nil {
		return &req, err
	}
	data := protocol.Data{Header: topic.Event.Key.Header, Key: topic.Event.Key.Array, Value: value}
	req.Data = &data
	return &req, nil
}

func (p *ProtoTopicFactoryObj) Topic(data []byte) (*protocol.Topic, error) {
	var tp protocol.Topic
	err := proto.Unmarshal(data, &tp)
	return &tp, err
}

func (p *ProtoTopicFactoryObj) Message(topic *protocol.Topic) (any, error) {
	m := p.Mt()
	err := anypb.UnmarshalTo(topic.Event.Message, m, proto.UnmarshalOptions{})
	if err != nil {
		return m, err
	}
	return m, nil
}
