package clustering

import (
	context "context"
	"fmt"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *DataServiceProvider) runPublish(topic *protocol.Topic) (*protocol.Response, error) {
	rc := make(chan *protocol.Response, 1)
	defer close(rc)

	tpf, registered := bootstrap.TopicFactoryRegistry[topic.Name]
	if !registered {
		return &protocol.Response{Successful: false}, fmt.Errorf("event factory not registered")
	}
	req, err := tpf().Request(topic)
	if err != nil {
		rc <- &protocol.Response{Successful: false}
	}
	resp, err := c.runCreate(req)
	if !resp.Successful {
		core.AppLog.Warn().Msgf("cannot save topic %v", resp)
		return resp, fmt.Errorf("cannot save topic")
	}
	rq := make(chan []core.Subscription, 3)
	defer close(rq)
	c.DRequest <- TopicRequest{Opt: TOPIC_REGISTER, Subs: rq, NodeId: topic.NodeId, Tag: topic.Tag, Name: topic.Name}
	subs := <-rq
	for _, sub := range subs {
		c.clientPublish(&sub, topic, rc)
		resp = <-rc
	}
	return &protocol.Response{Successful: true, Message: "topic delivered"}, nil
}
func (m *DataServiceProvider) clientPublish(target *core.Subscription, request *protocol.Topic, async chan *protocol.Response) {
	tcp, err := grpc.NewClient(target.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		async <- &protocol.Response{Successful: false, Message: err.Error()}
		return
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	resp, err := dsp.Send(context.Background(), request)
	async <- resp
}
