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

type QScore struct{
	TournamentId int64 `json:"TournamentId,string"`
	InstanceId   int64 `json:"InstanceId,string"`
	QWithTag
}

func (q *QScore) QCriteria(buff core.DataBuffer) error {
	buff.WriteString(q.Tag)
	buff.WriteString(T_SCORE_TAG)
	if q.TournamentId > 0 {
		buff.WriteInt64(q.TournamentId)
	}
	if q.InstanceId > 0 {
		buff.WriteInt64(q.InstanceId)
	}
	return nil
}
