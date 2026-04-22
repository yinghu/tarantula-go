package persistence

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func NewTransactionObjectFactory() *TransactionObjectFactory {
	mf := TransactionObjectFactory{}

	mf.Mo = func() proto.Message { return &protocol.Transaction{} }

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

func (p *TransactionObjectFactory) FromTransactionObject(trans *protocol.Transaction) (*protocol.KeyValue, error) {
	kv := protocol.KeyValue{}
	if trans.Meta.Id == 0 {
		return &kv, fmt.Errorf("none zeor id required")
	}
	buffer := core.NewBuffer(8)
	buffer.WriteUInt32(uint32(trans.Meta.Id))
	buffer.Flip()
	k, err := buffer.Read(0)
	if err != nil {
		return &kv, err
	}
	kv.Key = &protocol.Key{Array: k, Header: &protocol.Header{FactoryId: core.OBJECT_FACTORY_ID, ClassId: TRANSACTION_OBJECT_ID, Mutable: true}}
	obj, err := anypb.New(trans)
	if err != nil {
		return &kv, err
	}
	kv.Message = obj
	return &kv, nil
}
