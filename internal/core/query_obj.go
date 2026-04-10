package core

import "time"

type QueryObj struct {
	Id        uint32    `json:"-"`
	FactoryId uint32    `json:"-"`
	ClassId   uint32    `json:"-"`
	NodeId    string    `json:"NodeId"`
	Tag       string    `json:"Tag"`
	Topic     string    `json:"Topic"`
	Limit     int32     `json:"Limit"`
	Offset    int32     `json:"Offset"`
	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`

	Ee Event `json:"-"`
}

func (q *QueryObj) QRead(buff DataBuffer) error {
	id, err := buff.ReadUInt32()
	if err != nil {
		return err
	}
	q.Id = id
	fid, err := buff.ReadUInt32()
	if err != nil {
		return err
	}
	q.FactoryId = fid
	cid, err := buff.ReadUInt32()
	if err != nil {
		return err
	}
	q.ClassId = cid
	nid, err := buff.ReadString()
	if err != nil {
		return err
	}
	q.NodeId = nid
	tag, err := buff.ReadString()
	if err != nil {
		return err
	}
	q.Tag = tag
	topic, err := buff.ReadString()
	if err != nil {
		return err
	}
	q.Topic = topic
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
	if err := buff.WriteUInt32(q.Id); err != nil {
		return nil
	}
	if err := buff.WriteUInt32(q.FactoryId); err != nil {
		return nil
	}
	if err := buff.WriteUInt32(q.ClassId); err != nil {
		return nil
	}
	if err := buff.WriteString(q.NodeId); err != nil {
		return err
	}
	if err := buff.WriteString(q.Tag); err != nil {
		return err
	}
	if err := buff.WriteString(q.Topic); err != nil {
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

func (q *QueryObj) QNodeId() string {
	return q.NodeId
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

func (q *QueryObj) QEvent() Event {
	return q.Ee
}

func (q *QueryObj) QFilter(k, v []byte) bool {
	return true
}
