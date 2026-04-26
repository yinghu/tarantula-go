package clustering

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type JobResource struct {
	resource    *protocol.Job
	joining     map[uint64]*TransactionResource
	joinParties int
	confirmed   int
}

func (j *JobResource) join(meta *protocol.Meta) bool {
	t := j.joining[meta.Id]
	core.AppLog.Debug().Msgf("meta ID : %d STATE : %d CONFIRMED : %d FINISHED %d", meta.Id, meta.State, t.confirmed, t.finished)
	switch meta.State {
	case protocol.TCC_CONFIRMED:
		if t.confirmed > 0 {
			return false
		}
		t.resource.Meta.State = protocol.TCC_CONFIRMED
		t.confirmed++

	case protocol.TCC_FINISHED:
		if t.finished > 0 {
			return false
		}
		t.resource.Meta.State = protocol.TCC_FINISHED
		t.finished++
	default:
		return false
	}
	j.confirmed++
	return j.joinParties == j.confirmed
}
