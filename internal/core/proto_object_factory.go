package core

import (
	"gameclustering.com/internal/protocol"
)


type ProtoObjectFactory interface {
	//CompositeKey
	Request(obj *protocol.KeyValue) (*protocol.Request, error)
	Object(data []byte) (*protocol.KeyValue, error)
	Message(obj *protocol.KeyValue) (any, error)
	QueryFactory
}
