package core

import (
	"time"
)

const (
	EVENT_FACTORY_ID uint32 = 1
)

type EventListener interface {
	OnEvent(e Event)
	OnError(e Event, err error)
}

type EventService interface {
	EventCreator
	VerifyTicket(ticket string) (OnSession, error)
	EventListener
}

type Event interface {
	Inbound(buff DataBuffer) error
	Outbound(buff DataBuffer) error
	Persistentable
	OnListener(el EventListener)
	Listener() EventListener

	OnNodeId(t string)
	NodeId() string

	OnTag(t string)
	Tag() string

	OnTopic(t string)
	Topic() string

	OId() int64
	OnOId(id int64)
	OnRecipientId(id int64)
	RecipientId() int64
}

type EventCreator interface {
	Create(classId uint32, topic string) (Event, error)
}

type Pusher interface {
	Push(e Event)
}

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

type EventObj struct {
	Callback EventListener `json:"-"`
	PersistentableObj
	ENodeId string `json:"nodeId"`
	ETag    string `json:"tag"`
	ETopic  string `json:"topic"`
	EOid    int64  `json:"oid,string"`
}

func (s *EventObj) OnNodeId(t string) {
	s.ENodeId = t
}

func (s *EventObj) NodeId() string {
	return s.ENodeId
}

func (s *EventObj) OnTag(t string) {
	s.ETag = t
}

func (s *EventObj) Tag() string {
	return s.ETag
}

func (s *EventObj) OnTopic(t string) {
	s.ETopic = t
}

func (s *EventObj) Topic() string {
	return s.ETopic
}

func (s *EventObj) Inbound(buff DataBuffer) error {
	return nil
}
func (s *EventObj) Outbound(buff DataBuffer) error {
	return nil
}
func (s *EventObj) OnListener(el EventListener) {
	s.Callback = el
}
func (s *EventObj) Listener() EventListener {
	return s.Callback
}

func (s *EventObj) OnOId(oid int64) {
	s.EOid = oid
}

func (s *EventObj) OId() int64 {
	return s.EOid
}

func (s *EventObj) RecipientId() int64 {
	return 0
}

func (s *EventObj) OnRecipientId(recipientId int64) {

}

func (s *EventObj) FactoryId() uint32 {
	return EVENT_FACTORY_ID
}

type QueryObj struct {
	Id        uint32     `json:"-"`
	FactoryId uint32     `json:"-"`
	ClassId   uint32     `json:"-"`
	Tag       string     `json:"Tag"`
	Topic     string     `json:"Topic"`
	Limit     int32      `json:"Limit"`
	Offset    int32      `json:"Offset"`
	StartTime time.Time  `json:"StartTime"`
	EndTime   time.Time  `json:"EndTime"`
	Cc        chan Chunk `json:"-"`
	Ee        Event      `json:"-"`
}

func (q *QueryObj) QRead(buff DataBuffer) error {
	tag, err := buff.ReadString()
	if err != nil {
		return err
	}
	q.Tag = tag
	st, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	q.StartTime = time.UnixMilli(st)
	et, err := buff.ReadInt64()
	if err != nil {
		return err
	}
	q.EndTime = time.UnixMilli(et)
	lm, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	q.Limit = lm
	off, err := buff.ReadInt32()
	if err != nil {
		return err
	}
	q.Offset = off

	return nil
}

func (q *QueryObj) QWrite(buff DataBuffer) error {
	if err := buff.WriteString(q.Tag); err != nil {
		return err
	}
	if err := buff.WriteInt64(q.StartTime.UnixMilli()); err != nil {
		return err
	}
	if err := buff.WriteInt64(q.EndTime.UnixMilli()); err != nil {
		return err
	}
	if err := buff.WriteInt32(q.Limit); err != nil {
		return err
	}
	if err := buff.WriteInt32(q.Offset); err != nil {
		return err
	}
	return nil
}

func (q *QueryObj) QId() uint32 {
	return q.Id
}
func (q *QueryObj) QFactoryId() uint32 {
	return q.FactoryId
}
func (q *QueryObj) QClassId() uint32 {
	return q.ClassId
}
func (q *QueryObj) QTag() string {
	return q.Tag
}
func (q *QueryObj) QTopic() string {
	return q.Topic
}
func (q *QueryObj) QStartTime() time.Time {
	return q.StartTime
}
func (q *QueryObj) QEndTime() time.Time {
	return q.EndTime
}
func (q *QueryObj) QLimit() int32 {
	return q.Limit
}

func (q *QueryObj) QOffset() int32 {
	return q.Offset
}

func (q *QueryObj) QCc() chan Chunk {
	return q.Cc
}

func (q *QueryObj) QEvent() Event {
	return q.Ee
}

func (q *QueryObj) QFilter(k, v []byte) bool {
	return true
}
