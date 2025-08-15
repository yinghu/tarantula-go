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
	buff.WriteInt64(q.TournamentId)
	buff.WriteInt64(q.InstanceId)
	buff.WriteInt64(q.SystemId)
	return nil
}
