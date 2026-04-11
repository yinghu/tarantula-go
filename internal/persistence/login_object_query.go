package persistence

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type LoginObjectQuery struct {
	core.QueryObj
}

func (q *LoginObjectQuery) QFilter(k, v []byte) bool {
	mf := NewLoginObjectFactory()
	t, err := mf.Object(v)
	if err != nil {
		core.AppLog.Warn().Msgf("wrong decode format %s", err)
		return false
	}
	obj, err := mf.Message(t)
	if err != nil {
		core.AppLog.Warn().Msgf("wrong decode format %s", err)
		return false
	}
	m, ok := obj.(*protocol.LoginObject)
	if !ok {
		core.AppLog.Warn().Msgf("wrong login object format %s", err)
		return false
	}
	core.AppLog.Debug().Msgf("filter here %v", m)
	return true
}

func (q *LoginObjectQuery) QList(list core.List) error {
	if err := q.QueryObj.QList(list); err != nil {
		return err
	}
	mf := NewLoginObjectFactory()
	for _, data := range q.Payload.Data.List {
		kv, err := mf.Object(data.Value)
		if err != nil {
			continue
		}
		core.AppLog.Debug().Msgf("key value : %v", kv)
		e, err := mf.Message(kv)
		if err != nil {
			continue
		}
		_, ok := e.(*protocol.LoginObject)
		if !ok {
			core.AppLog.Warn().Msg("casting error")
		}
		if !list(kv.Key.Header, e) {
			break
		}
	}
	return nil
}
