package persistence

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func NewTransactionObjectFactory() *TransactionObjectFactory {
	mf := TransactionObjectFactory{}

	mf.Mo = func() proto.Message { return &protocol.Meta{} }

	mq := TransactionObjectQuery{}
	mq.Id = TRANSACTION_OBJECT_FACTORY_NAME
	mq.FactoryId = core.OBJECT_FACTORY_ID
	mq.ClassId = TRANSACTION_OBJECT_ID
	mq.Topic = TRANSACTION_OBJECT_FACTORY_NAME
	mf.Q = &mq
	return &mf
}

type TransactionObjectFactory struct {
	ProtoObjectFactoryObj
}

func (p *TransactionObjectFactory) FromId(tid uint64) (*protocol.Request, error) {
	req := protocol.Request{Opt: core.GET_DATA_REQUEST}
	buffer := core.NewBuffer(8)
	buffer.WriteUInt64(tid)
	buffer.Flip()
	k, err := buffer.Read(0)
	if err != nil {
		return &req, err
	}
	req.Prefix = util.Hash(k)
	req.Data = &protocol.Data{Key: k, Header: &protocol.Header{FactoryId: core.OBJECT_FACTORY_ID, ClassId: TRANSACTION_OBJECT_ID, Mutable: true}}
	return &req, nil
}
func (p *TransactionObjectFactory) FromTransactionObject(trans *protocol.Transaction) (*protocol.KeyValue, error) {
	kv := protocol.KeyValue{}
	if trans.Meta.Id == 0 {
		return &kv, fmt.Errorf("none zeor id required")
	}
	buffer := core.NewBuffer(8)
	buffer.WriteUInt64(trans.Meta.Id)
	buffer.Flip()
	k, err := buffer.Read(0)
	if err != nil {
		return &kv, err
	}
	kv.Key = &protocol.Key{Array: k, Header: &protocol.Header{FactoryId: core.OBJECT_FACTORY_ID, ClassId: TRANSACTION_OBJECT_ID, Mutable: true}}
	obj, err := anypb.New(trans.Meta)
	if err != nil {
		return &kv, err
	}
	kv.Message = obj
	return &kv, nil
}
