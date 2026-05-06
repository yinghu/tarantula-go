package persistence

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type CommodityObjectQuery struct {
	ObjectQueryObj
}

func (q *CommodityObjectQuery) QFilter(k, v []byte) bool {
	mf := NewCommodityObjectFactory()
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
	m, ok := obj.(*protocol.Commodity)
	if !ok {
		core.AppLog.Warn().Msgf("wrong login object format %s", err)
		return false
	}
	core.AppLog.Debug().Msgf("filter here %v", m)
	return true
}
