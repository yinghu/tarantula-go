package event

import (
	"time"

	"gameclustering.com/internal/core"
)

const (
	TAG_MESSAGE_QID    int32 = 0
	TAG_TOURNAMENT_QID int32 = 1
	Q_TOURNAMENT_QID   int32 = 2
)

func CreateQuery(qid int32) Query {
	switch qid {
	case TAG_MESSAGE_QID:
		q := QWithTag{Id: qid, Tag: MESSAGE_ETAG, Cc: make(chan Chunk, 3)}
		return &q
	case TAG_TOURNAMENT_QID:
		q := QWithTag{Id: qid, Tag: TOURNAMENT_ETAG, Cc: make(chan Chunk, 3)}
		return &q
	case Q_TOURNAMENT_QID:
		q := QTournament{}
		q.Id = qid
		q.Tag = TOURNAMENT_ETAG
		q.Cc = make(chan Chunk, 3)
		return &q
	default:
		q := QWithTag{Id: qid, Tag: MESSAGE_ETAG, Cc: make(chan Chunk, 3)}
		return &q
	}
}

type Query interface {
	QId() int32
	QTag() string
	QStartTime() time.Time
	QEndTime() time.Time
	QLimit() int32
	QCriteria(b core.DataBuffer) error
	QCc() chan Chunk
}

type QWithTag struct {
	Id        int32      `json:"-"`
	Tag       string     `json:"Tag"`
	Limit     int32      `json:"Limit"`
	StartTime time.Time  `json:"StartTime"`
	EndTime   time.Time  `json:"EndTime"`
	Cc        chan Chunk `json:"-"`
}

func (q *QWithTag) QCriteria(buff core.DataBuffer) error {
	buff.WriteString(q.Tag)
	return nil
}

func (q *QWithTag) QId() int32 {
	return q.Id
}

func (q *QWithTag) QTag() string {
	return q.Tag
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

func (q *QWithTag) QCc() chan Chunk {
	return q.Cc
}
