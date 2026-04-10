package core

import (
	"gameclustering.com/internal/protocol"
)

const (
	EVENT_FACTORY_ID uint32 = 1
)

type ProtoTopicFactory interface {
	CompositeKey
	Request(topic *protocol.Topic) (*protocol.Request, error)
	Topic(data []byte) (*protocol.Topic, error)
	Message(topic *protocol.Topic) (any, error)
	QueryFactory
}
