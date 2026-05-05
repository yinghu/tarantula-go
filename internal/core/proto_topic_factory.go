package core

import (
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type ProtoTopicFactory interface {
	//CompositeKey
	Request(topic *protocol.Topic) (*protocol.Request, error)
	Hash(h MessageHash) uint32
	Topic(data []byte) (*protocol.Topic, error)
	Message(topic *protocol.Topic) (proto.Message, error)
	QueryFactory
}
