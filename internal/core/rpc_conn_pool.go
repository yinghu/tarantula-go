package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	MIN_POOL_SIZE       int = 1
	MAX_POOL_SIZE       int = 1
	RPC_CONNECT_RETRIES int = 3
)

type RpcConn struct {
	Conn *grpc.ClientConn
	Seq  int
}

type RpcConnPool struct {
	Target  string
	Tag     string
	NodeId  string
	MinSize int
	MaxSize int
	index   int
	pool    map[string]*RpcConn
	sync.RWMutex
}

func (p *RpcConnPool) connect(target string) (*grpc.ClientConn, error) {
	retries := RPC_CONNECT_RETRIES
	creds, err := credentials.NewClientTLSFromFile(CERT_NAME, "")
	if err != nil {
		panic(err.Error())
	}
	for {
		tcp, err := grpc.NewClient(target, grpc.WithTransportCredentials(creds), grpc.WithUnaryInterceptor(p.audit))
		if err != nil {
			retries--
			if retries > 0 {
				AppLog.Warn().Msgf("retrying to connect gprc %s with retried times %d", err.Error(), retries)
				time.Sleep(3 * time.Second)
				continue
			}
			return tcp, err
		}
		return tcp, nil
	}
}

func (p *RpcConnPool) Start() {
	p.Lock()
	defer p.Unlock()
	if p.MinSize <= 0 {
		p.MinSize = MIN_POOL_SIZE
	}
	if p.MaxSize <= 0 {
		p.MaxSize = MAX_POOL_SIZE
	}
	p.pool = make(map[string]*RpcConn)
	p.index = 0
}

func (p *RpcConnPool) Conn() (*RpcConn, error) {
	p.Lock()
	defer p.Unlock()
	ckey := fmt.Sprintf("%s_%d", p.Target, p.index)
	c, created := p.pool[ckey]
	if !created {
		cx, err := p.connect(p.Target)
		if err != nil {
			return c, err
		}
		c = &RpcConn{Conn: cx, Seq: p.index}
		p.pool[ckey] = c
	}
	return c, nil
}
func (p *RpcConnPool) Release() {
	p.Lock()
	defer p.Unlock()
	for _, c := range p.pool {
		dsp := protocol.NewPostofficeServiceClient(c.Conn)
		resp, err := dsp.Disconnect(context.Background(), &protocol.Topic{Tag: p.Tag, NodeId: p.NodeId})
		if err != nil {
			continue
		}
		AppLog.Debug().Msgf("disconnecting topic %v", resp)
		c.Conn.Close()
	}
	clear(p.pool)
}

func (p *RpcConnPool) audit(ctx context.Context, method string, req, replay any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//AppLog.Debug().Msgf("call before :%s",method)
	fmt.Printf("ssmethod 111>%s", method)
	err := invoker(ctx, method, req, replay, cc, opts...)
	fmt.Printf("ssmethod 222>%s", method)
	return err
}
