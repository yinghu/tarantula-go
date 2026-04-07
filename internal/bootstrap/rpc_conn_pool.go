package bootstrap

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	MIN_POOL_SIZE int = 1
	MAX_POOL_SIZE int = 1
)

type RpcConn struct {
	Conn *grpc.ClientConn
	Seq  int
}

type RpcConnPool struct {
	Target  string
	MinSize int
	MaxSize int
	index   int
	pool    map[string]*RpcConn
	sync.RWMutex
}

func (p *RpcConnPool) connect(target string) (*grpc.ClientConn, error) {
	retries := RPC_CONNECT_RETRIES
	for {
		tcp, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			retries--
			if retries > 0 {
				core.AppLog.Warn().Msgf("retrying to connect gprc %s with retried times %d", err.Error(), retries)
				time.Sleep(3 * time.Second)
				continue
			}
			return tcp, err
		}
		break
	}
	return nil, fmt.Errorf("should not be here")
}

func (p *RpcConnPool) Start() error {
	p.Lock()
	defer p.Unlock()
	if p.MinSize <= 0 {
		p.MinSize = MIN_POOL_SIZE
	}
	if p.MaxSize <= 0 {
		p.MaxSize = MAX_POOL_SIZE
	}
	p.pool = make(map[string]*RpcConn)
	for i := range p.MinSize {
		cx, err := p.connect(p.Target)
		if err != nil {
			return err
		}
		ckey := fmt.Sprintf("%s_%d", p.Target, i)
		p.pool[ckey] = &RpcConn{Conn: cx, Seq: i}
	}
	p.index = 0
	return nil
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
func (p *RpcConnPool) Release(t *protocol.Topic) {
	p.Lock()
	defer p.Unlock()
	for _, c := range p.pool {
		dsp := protocol.NewPostofficeServiceClient(c.Conn)
		resp, err := dsp.Disconnect(context.Background(), t)
		if err != nil {
			continue
		}
		core.AppLog.Debug().Msgf("disconnecting topic %v", resp)
		c.Conn.Close()
	}
	clear(p.pool)
}
