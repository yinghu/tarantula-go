package core

import (
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type ProtoTopicFactory interface {
	CompositeKey
	Request(topic *protocol.Topic) (*protocol.Request, error)
	Topic(data []byte) (*protocol.Topic, error)
	Query(criteria []byte) (Query, error)
}

type ProtoTopicFactoryObj struct {
	Target *protocol.Topic
}

func (p *ProtoTopicFactoryObj) WriteKey(key DataBuffer) error {
	return nil
}

func (p *ProtoTopicFactoryObj) ReadKey(key DataBuffer) error {
	return nil
}

func (p *ProtoTopicFactoryObj) Topic(data []byte) (*protocol.Topic, error) {
	var tp protocol.Topic
	err := proto.Unmarshal(data, &tp)
	return &tp, err
}
