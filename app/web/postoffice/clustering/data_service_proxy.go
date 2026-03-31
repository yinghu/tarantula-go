package clustering

import (
	context "context"
	"fmt"
	"io"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *DataServiceProvider) runCreate(set *protocol.Request) (*protocol.Response, error) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	var rt uint32
	if set.Prefix > 0 {
		rt = set.Prefix
		core.AppLog.Debug().Msgf("using prefix %d", set.Prefix)
	} else {
		rt = c.Mll.RingToken(set.Data.Key)
	}
	ch := make(chan *protocol.Response, 3)
	defer close(ch)
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: rt, Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		c.clientCreate(&ringNode, set, ch)
		resp := <-ch
		if !resp.Successful {
			retry.Err = resp.Message
			retry.Reties--
			continue
		}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientCreate(&slave, set, ch)
			resp = <-ch
			if !resp.Successful {
				core.AppLog.Debug().Msgf("error on slave %s", resp.Message)
			}
		}
		break
	}
	core.AppLog.Printf("retry %s, %d %v", retry.Err, retry.Reties, retry.Suc)
	return &protocol.Response{Successful: retry.Suc, Message: retry.Err}, nil
}

func (m *DataServiceProvider) clientCreate(target *core.Node, request *protocol.Request, ch chan *protocol.Response) {
	task := Task{target: target.RpcEndpoint, execute: func(tcp *grpc.ClientConn, opt int) error {
		dsp := protocol.NewDataServiceClient(tcp)
		resp, _ := dsp.Create(context.Background(), request)
		ch <- resp
		return nil
	}}
	m.WTask <- task
}

func (c *DataServiceProvider) runUpdate(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientUpdate(&ringNode, set)
		if err != nil {
			retry.Err = err.Error()
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientUpdate(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err, retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err}
}

func (m *DataServiceProvider) clientUpdate(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Update(context.Background(), request)
}

func (c *DataServiceProvider) runDelete(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientDelete(&ringNode, set)
		if err != nil {
			retry.Err = err.Error()
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientDelete(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err, retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err}
}

func (m *DataServiceProvider) clientDelete(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Delete(context.Background(), request)
}

func (c *DataServiceProvider) runReset(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: c.Mll.RingToken(set.Data.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		resp, err := c.clientReset(&ringNode, set)
		if err != nil {
			retry.Err = err.Error()
			retry.Reties--
			continue
		}
		ch <- resp
		retry.Suc = true
		if !resp.Successful {
			break
		}
		slaves := nodes[1:]
		for _, slave := range slaves {
			c.clientReset(&slave, set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err, retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err}
}

func (m *DataServiceProvider) clientReset(target *core.Node, request *protocol.Request) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Reset(context.Background(), request)
}

func (c *DataServiceProvider) runGet(set *protocol.Request, ch chan *protocol.Response) {
	rq := make(chan []core.Node, 3)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}
	for retry.Reties > 0 {
		kh := c.Mll.RingToken(set.Data.Key)
		c.Mll.MRequest <- core.RingRequest{Opt: REPLICA_RING_OPT, Token: kh, Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		tcp, err := grpc.NewClient(ringNode.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			retry.Reties++
			continue
		}
		defer tcp.Close()
		dsp := protocol.NewDataServiceClient(tcp)
		stream, err := dsp.Get(context.Background(), set)
		if err != nil {
			retry.Reties++
			continue
		}
		crt := protocol.Response{Successful: false}
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				core.AppLog.Debug().Msgf("run get streaming error %s", err.Error())
				crt.Code = 400001
				crt.Message = err.Error()
				break
			}
			ch <- data
			if !data.Successful {
				break
			}
		}
		ch <- &crt
		retry.Suc = true
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err, retry.Reties)
	ch <- &protocol.Response{Successful: false, Message: retry.Err}
}

func (c *DataServiceProvider) runPull(target string, set *protocol.Request, ch chan *protocol.Response) {
	core.AppLog.Debug().Msgf("run remote pull %s >= %d < %d", target, set.Prefix, set.Opt)
	tcp, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		ch <- &protocol.Response{Successful: false, Message: err.Error()}
		return
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	stream, err := dsp.Pull(context.Background(), set)
	if err != nil {
		ch <- &protocol.Response{Successful: false, Message: err.Error()}
		return
	}
	crt := protocol.Response{Successful: false}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("run pull streaming error %s", err.Error())
			crt.Message = err.Error()
			break
		}
		ch <- data
		if !data.Successful {
			break
		}
	}
	ch <- &crt
}

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
