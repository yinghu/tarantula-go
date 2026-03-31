package bootstrap

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type ProtoTopicFactory struct {
	Cluster core.ClusterService
}

func (p *ProtoTopicFactory) ToRequest(topic *protocol.Topic) (*protocol.Request, error) {
	req := protocol.Request{Opt: core.CREATE_DATA_REQUEST, Prefix: p.Cluster.RingToken([]byte(topic.Name))}
	if topic.Event.Id <= 0 {
		return &req, fmt.Errorf("id cannot be less than 0")
	}
	buff := core.NewBuffer(8)
	if err := buff.WriteUInt64(topic.Event.Id); err != nil {
		return &req, err
	}
	buff.Flip()
	key, err := buff.Read(0)
	if err != nil {
		return &req, err
	}
	value, err := proto.Marshal(topic)
	if err != nil {
		return &req, err
	}
	data := protocol.Data{Header: topic.Event.Header, Key: key, Value: value}
	req.Data = &data
	return &req, nil
}

func (p *ProtoTopicFactory) FromData(data []byte) (*protocol.Topic, error) {
	tp := protocol.Topic{}
	err := proto.Unmarshal(data, &tp)
	return &tp, err
}
