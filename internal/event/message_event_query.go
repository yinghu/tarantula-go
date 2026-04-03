package event

import (
	"gameclustering.com/internal/core"
)

type MessageEventQuery struct {
	core.QueryObj
}

func (q *MessageEventQuery) QFilter(k, v []byte) bool {

	core.AppLog.Debug().Msgf("filter here")
	return true
}
