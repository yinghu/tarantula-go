package persistence

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"
)

type ProtoObjectFactoryObj struct {
}

func (p *ProtoObjectFactoryObj) Request(obj *protocol.KeyValue) (*protocol.Request, error) {
	req := protocol.Request{Opt: core.CREATE_OBJ_REQUEST}
	req.Prefix = util.Hash(obj.Key.Array)
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

func (p *ProtoObjectFactoryObj) Export(query core.Query) ([]byte, error) {
	return nil, nil
	//return Export(query, core.COMPOSIT_KEY_MAX)
}
