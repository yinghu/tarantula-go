package core

import "time"

type Query interface {
	QId() uint32
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
	QEvent() Event
	QFilter(k, v []byte) bool
}

type QueryFactory interface {
	Export(query Query) ([]byte, error)
	Import(criteria []byte) (Query, error)
	Query() Query
}
