package core

import (
	"time"

	"gameclustering.com/internal/protocol"
)

type List func(h *protocol.Header, m any) bool

type Query interface {
	QId() string
	QFactoryId() uint32
	QClassId() uint32
	QNodeId() string
	QTag() string
	QTopic() string
	QStartTime() time.Time
	QEndTime() time.Time
	QLimit() int32
	QOffset() int32
	QRead(b DataBuffer) error
	QWrite(b DataBuffer) error
	QFilter(k, v []byte) bool

	//read
	QList(list List) error
	QResponse(resp *protocol.Response)
}

type QueryFactory interface {
	Export(query Query) ([]byte, error)
	Import(criteria []byte) (Query, error)
	Query() Query
	Set(resp *protocol.Response) Query
}
