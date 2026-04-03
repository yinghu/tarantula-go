package event

import (
	"gameclustering.com/internal/core"
)

type MessageEventQuery struct {
	core.QueryObj
}

func (q *MessageEventQuery) QFilter(k, v []byte) bool {
	mf := MessageEventFactory{}
	t, err := mf.Topic(v)
	if err != nil {
		core.AppLog.Warn().Msgf("wrong decode format %s", err)
		return false
	}
	core.AppLog.Debug().Msgf("filter here %v", t)
	return true
}
