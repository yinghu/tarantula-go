package event

import (
	"time"

	"gameclustering.com/internal/core"
)

const (
	TAG_MESSAGE    int32 = 0
	TAG_TOURNAMENT int32 = 1
)

func CreateQuery(qid int32) Query {
	switch qid {
	case TAG_MESSAGE:
		q := QueryObj{Id: qid, Tag: MESSAGE_ETAG, Cc: make(chan Chunk, 3)}
		return &q
	case TAG_TOURNAMENT:
		q := QueryObj{Id: qid, Tag: TOURNAMENT_ETAG, Cc: make(chan Chunk, 3)}
		return &q
	default:
		q := QueryObj{Id: qid, Tag: MESSAGE_ETAG, Cc: make(chan Chunk, 3)}
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

type QueryObj struct {
	Id        int32      `json:"-"`
	Tag       string     `json:"Tag"`
	Limit     int32      `json:"Limit"`
	StartTime time.Time  `json:"StartTime"`
	EndTime   time.Time  `json:"EndTime"`
	Cc        chan Chunk `json:"-"`
}

func (q *QueryObj) QCriteria(buff core.DataBuffer) error {
	buff.WriteString(q.Tag)
	return nil
}

func (q *QueryObj) QId() int32 {
	return q.Id
}

func (q *QueryObj) QTag() string {
	return q.Tag
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

func (q *QueryObj) QCc() chan Chunk {
	return q.Cc
}
