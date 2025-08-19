package event

import "gameclustering.com/internal/core"

type QTournament struct {
	TournamentId int64 `json:"TournamentId,string"`
	InstanceId   int64 `json:"InstanceId,string"`
	SystemId     int64 `json:"SystemId,string"`
	QWithTag
}

func (q *QTournament) QCriteria(buff core.DataBuffer) error {
	buff.WriteString(q.Tag)
	if q.TournamentId > 0 {
		buff.WriteInt64(q.TournamentId)
	}
	if q.InstanceId > 0 {
		buff.WriteInt64(q.InstanceId)
	}
	if q.SystemId > 0 {
		buff.WriteInt64(q.SystemId)
	}
	return nil
}

type QJoin struct{
	SystemId     int64 `json:"SystemId,string"`
	QWithTag
}

func (q *QJoin) QCriteria(buff core.DataBuffer) error {
	buff.WriteString(q.Tag)
	buff.WriteString(T_JOIN_TAG)
	if q.SystemId > 0{
		buff.WriteInt64(q.SystemId)
	}
	return nil
}
