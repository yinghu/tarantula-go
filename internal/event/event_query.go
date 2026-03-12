package event

import (
	"time"

	"gameclustering.com/internal/core"
)

const (
	//query with tag
	TAG_MESSAGE_QID    uint32 = 0
	TAG_LOGIN_QID      uint32 = 1
	TAG_TOURNAMENT_QID uint32 = 2
	TAG_INVENTORY_QID  uint32 = 3
	TAG_KICKOFF_QID    uint32 = 4
	TAG_REGISTER_QID   uint32 = 5

	//query criteria
	Q_TOURNAMENT_QID uint32 = 100
	QT_SCORE_QID     uint32 = 101
)

func Import(q core.Query, data []byte, bufferiSize int) error {
	buff := core.NewBuffer(bufferiSize)
	if err := buff.Write(data); err != nil {
		return err
	}
	buff.Flip()
	return q.QRead(buff)
}
func Export(q core.Query, buffSize int) ([]byte, error) {
	var v []byte
	buff := core.NewBuffer(buffSize)
	if err := q.QWrite(buff); err != nil {
		return v, nil
	}
	buff.Flip()
	return buff.Read(0)
}

func CreateQuery(qid uint32) core.Query {
	switch qid {
	case TAG_MESSAGE_QID:
		q := QWithTag{Id: qid, Tag: MESSAGE_ETAG, Cc: make(chan core.Chunk, 3), FactoryId: core.EVENT_FACTORY_ID, ClassId: MESSAGE_CID}
		return &q
	case TAG_LOGIN_QID:
		q := QWithTag{Id: qid, Tag: LOGIN_ETAG, Cc: make(chan core.Chunk, 3)}
		return &q

	case TAG_TOURNAMENT_QID:
		q := QWithTag{Id: qid, Tag: TOURNAMENT_ETAG, Cc: make(chan core.Chunk, 3)}
		return &q
	case TAG_INVENTORY_QID:
		q := QWithTag{Id: qid, Tag: INVENTORY_ETAG, Cc: make(chan core.Chunk, 3)}
		return &q
	case TAG_KICKOFF_QID:
		q := QWithTag{Id: qid, Tag: KICKOFF_ETAG, Cc: make(chan core.Chunk, 3)}
		return &q
	case TAG_REGISTER_QID:
		q := QWithTag{Id: qid, Tag: REGISTER_ETAG, Cc: make(chan core.Chunk, 3)}
		return &q
	case Q_TOURNAMENT_QID:
		q := QTournament{}
		q.Id = qid
		q.Tag = TOURNAMENT_ETAG
		q.Cc = make(chan core.Chunk, 3)
		return &q
	case QT_SCORE_QID:
		q := QScore{}
		q.Id = qid
		q.Tag = TOURNAMENT_ETAG
		q.Cc = make(chan core.Chunk, 3)
		return &q
	default:
		q := QWithTag{Id: qid, Tag: MESSAGE_ETAG, Cc: make(chan core.Chunk, 3)}
		return &q
	}
}

type QWithTag struct {
	Id        uint32          `json:"-"`
	FactoryId int32           `json:"-"`
	ClassId   int32           `json:"-"`
	Tag       string          `json:"Tag"`
	Topic     string          `json:"Topic"`
	Limit     int32           `json:"Limit"`
	Offset    int32           `json:"Offset"`
	StartTime time.Time       `json:"StartTime"`
	EndTime   time.Time       `json:"EndTime"`
	Cc        chan core.Chunk `json:"-"`
	Ee        core.Event      `json:"-"`
}

func (q *QWithTag) QRead(buff core.DataBuffer) error {
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

func (q *QWithTag) QWrite(buff core.DataBuffer) error {
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

func (q *QWithTag) QId() uint32 {
	return q.Id
}
func (q *QWithTag) QFactoryId() int32 {
	return q.FactoryId
}
func (q *QWithTag) QClassId() int32 {
	return q.ClassId
}
func (q *QWithTag) QTag() string {
	return q.Tag
}
func (q *QWithTag) QTopic() string {
	return q.Topic
}
func (q *QWithTag) QStartTime() time.Time {
	return q.StartTime
}
func (q *QWithTag) QEndTime() time.Time {
	return q.EndTime
}
func (q *QWithTag) QLimit() int32 {
	return q.Limit
}

func (q *QWithTag) QOffset() int32 {
	return q.Offset
}

func (q *QWithTag) QCc() chan core.Chunk {
	return q.Cc
}

func (q *QWithTag) QEvent() core.Event {
	return q.Ee
}

func (q *QWithTag) QFilter(k, v []byte) bool {
	buff := core.NewBuffer(100)
	buff.Write(k)
	buff.Flip()
	tag, _ := buff.ReadString()
	oid, _ := buff.ReadInt64()
	rev, _ := buff.ReadInt64()
	core.AppLog.Debug().Msgf("filter %s %d %d", tag, oid, rev)
	return tag == q.Tag
}
