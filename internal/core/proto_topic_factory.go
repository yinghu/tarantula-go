package core

import (
	"gameclustering.com/internal/protocol"
)

type ProtoTopicFactory interface {
	CompositeKey
	Request(topic *protocol.Topic) (*protocol.Request, error)
	Topic(data []byte) (*protocol.Topic, error)
	Export(query Query) ([]byte, error)
	Import(criteria []byte) (Query, error)
}
