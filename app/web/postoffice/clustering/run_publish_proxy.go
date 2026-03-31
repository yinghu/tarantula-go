package clustering

import (
	context "context"
	"fmt"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
)

func (c *DataServiceProvider) runPublish(topic *protocol.Topic) (*protocol.Response, error) {
	rc := make(chan *protocol.Response, 1)
	defer close(rc)

	tpf := bootstrap.ProtoTopicFactory{}
	req, err := tpf.ToRequest(topic)
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
		core.AppLog.Debug().Msgf("publish %v", resp)
	}
	return &protocol.Response{Successful: true, Message: "topic delivered"}, nil
}
func (m *DataServiceProvider) clientPublish(target *core.Subscription, request *protocol.Topic, async chan *protocol.Response) {
	task := Task{target: target.Endpoint, execute: func(tcp *grpc.ClientConn, opt int) error {
		if opt == NO_TCP_CONNECT {
			async <- &protocol.Response{Successful: false, Message: "no tcp"}
			return fmt.Errorf("no tcp from %s", target.Endpoint)
		}
		dsp := protocol.NewDataServiceClient(tcp)
		resp, err := dsp.Send(context.Background(), request)
		async <- resp
		core.AppLog.Debug().Msgf("SEND : %v", resp)
		return err
	}}
	m.WTask <- task
}
