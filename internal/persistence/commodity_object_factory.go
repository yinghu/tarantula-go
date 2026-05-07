package persistence

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

func NewCommodityObjectFactory() *CommodityObjectFactory {
	mf := CommodityObjectFactory{}

	mf.Mo = func() proto.Message { return &protocol.Commodity{} }

	mq := CommodityObjectQuery{}
	mq.Id = COMMODITY_OBJECT_FACTORY_NAME
	mq.FactoryId = core.OBJECT_FACTORY_ID
	mq.ClassId = COMMODITY_OBJECT_ID
	mq.Topic = COMMODITY_OBJECT_FACTORY_NAME
	mf.Q = &mq
	return &mf
}

type CommodityObjectFactory struct {
	ProtoObjectFactoryObj
}
