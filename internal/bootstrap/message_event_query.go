package bootstrap

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
)

type MessageEventQuery struct {
	event.QWithTag
}

func (q *MessageEventQuery) QFilter(k, v []byte) bool {
	core.AppLog.Debug().Msgf("filter here")
	return true
}
