package core

import "time"

type Query interface {
	QId() uint32
	QFactoryId() uint32
	QClassId() uint32
	QTag() string
	QTopic() string
	QStartTime() time.Time
	QEndTime() time.Time
	QLimit() int32
	QOffset() int32
	QRead(b DataBuffer) error
	QWrite(b DataBuffer) error
	QCc() chan Chunk
	QEvent() Event
	QFilter(k, v []byte) bool
}

type Postoffice interface {
	Send(e Event) error
	List(q Query)
	Recover(q Query)
	Load(e Query)
}
