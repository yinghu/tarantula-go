package event

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"
)

type ProtoTopicFactoryObj struct {
	Target *protocol.Topic
}

func (p *ProtoTopicFactoryObj) Request(topic *protocol.Topic) (*protocol.Request, error) {
	p.Target = topic
	req := protocol.Request{Opt: core.CREATE_DATA_REQUEST}
	if topic.Event.Id <= 0 {
		return &req, fmt.Errorf("id cannot be less than 0")
	}
	buff := core.NewBuffer(core.COMPOSIT_KEY_MAX)
	p.WriteKey(buff)
	buff.Flip()
	key, err := buff.Read(0)
	if err != nil {
		return &req, err
	}
	req.Prefix = util.Hash(key)
	value, err := proto.Marshal(topic)
	if err != nil {
		return &req, err
	}
	data := protocol.Data{Header: topic.Event.Header, Key: key, Value: value}
	req.Data = &data
	return &req, nil
}

func (p *ProtoTopicFactoryObj) WriteKey(key core.DataBuffer) error {
	return key.WriteUInt64(p.Target.Event.Id)
}

func (p *ProtoTopicFactoryObj) ReadKey(key core.DataBuffer) error {
	return nil
}

func (p *ProtoTopicFactoryObj) Topic(data []byte) (*protocol.Topic, error) {
	var tp protocol.Topic
	err := proto.Unmarshal(data, &tp)
	return &tp, err
}

func (p *ProtoTopicFactoryObj) Export(query core.Query) ([]byte, error) {
	return Export(query, core.COMPOSIT_KEY_MAX)
}
