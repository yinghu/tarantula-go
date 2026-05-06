package persistence

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	LOGIN_OBJECT_ID           uint32 = 1
	LOGIN_OBJECT_FACTORY_NAME        = "obj_login"
)

type MessageObject func() proto.Message

type ProtoObjectFactoryObj struct {
	Target *protocol.KeyValue
	core.QueryFactoryObj
	Mo MessageObject
}

func (p *ProtoObjectFactoryObj) Request(obj *protocol.KeyValue) (*protocol.Request, error) {
	p.Target = obj
	req := protocol.Request{Opt: core.CREATE_DATA_REQUEST}
	value, err := proto.Marshal(obj)
	if err != nil {
		return &req, err
	}
	data := protocol.Data{Header: obj.Key.Header, Key: obj.Key.Array, Value: value}
	req.Data = &data
	return &req, nil
}

func (p *ProtoObjectFactoryObj) Object(data []byte) (*protocol.KeyValue, error) {
	var kv protocol.KeyValue
	err := proto.Unmarshal(data, &kv)
	return &kv, err
}

func (p *ProtoObjectFactoryObj) Message(obj *protocol.KeyValue) (any, error) {
	m := p.Mo()
	err := anypb.UnmarshalTo(obj.Message, m, proto.UnmarshalOptions{})
	if err != nil {
		return m, err
	}
	return m, nil
}
func (p *ProtoObjectFactoryObj) Hash(mh core.MessageHash) uint32 {
	return mh.RingToken(p.Target.Key.Array)
}

func (p *ProtoObjectFactoryObj) FromMessage(m proto.Message, h *protocol.Header) (*protocol.KeyValue, error) {
	kv := protocol.KeyValue{}
	kv.Key = &protocol.Key{Header: h}
	obj, err := anypb.New(m)
	if err != nil {
		return &kv, err
	}
	kv.Message = obj
	return &kv, nil
}
