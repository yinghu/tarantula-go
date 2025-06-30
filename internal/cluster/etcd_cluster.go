package cluster

import (
	"context"
	"encoding/json"

	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gameclustering.com/internal/core"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type LocalNode struct {
	core.Node
	pingCount    *int8  `json:"-"`
	timeoutCount *uint8 `json:"-"`
}

type EtcdCluster struct {
	started  atomic.Bool
	starting *sync.WaitGroup
	ClusterManager
}

func newCluster(group string, etcEndpoints []string, local LocalNode, kl core.ClusterListener) core.Cluster {
	cmanager := ClusterManager{group: group, EtcdEndpoints: etcEndpoints, local: local}
	etc := EtcdCluster{ClusterManager: cmanager}
	etc.started.Store(false)
	etc.clistener = kl
	etc.lock = &sync.Mutex{}
	etc.cluster = make(map[string]LocalNode)
	etc.partition = make([]string, core.CLUSTER_PARTITION_NUM)
	etc.quit = make(chan bool)
	etc.starting = &sync.WaitGroup{}
	etc.starting.Add(1)
	return &etc
}

func (c *EtcdCluster) Join() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	nd, _ := json.Marshal(c.local)
	go func() {
		tik := time.NewTicker(1 * time.Second)
		defer tik.Stop()
		for r := range 7 {
			w := <-tik.C
			core.AppLog.Printf("Waiting for member joining [%s]: %v\n", c.Group(), w)
			if r == 0 {
				cli.Put(context.Background(), c.group+"#join", string(nd))
			} else {
				cli.Put(context.Background(), c.group+"#ping", string(nd))
			}
		}
		c.starting.Done()
	}()
	go func() {
		c.starting.Wait() //blocked
		tik := time.NewTicker(1 * time.Second)
		pct := 5
		for {
			select {
			case <-c.quit:
				cli.Close()
				return
			case p := <-tik.C:
				cli.Put(context.Background(), c.group+"#ping", c.local.Name)
				pct--
				if pct == 0 {
					pct = 5
					c.lock.Lock()
					for n := range c.cluster {
						if n != c.local.Name {
							cn := c.cluster[n]
							if *cn.timeoutCount == 3 {
								core.AppLog.Printf("Node timeout %d %d %s %v\n", *cn.pingCount, *cn.timeoutCount, cn.Name, p)
								delete(c.cluster, n)
								c.grouping()
								c.Listener().MemberLeft(cn.Node)
							} else {
								if *cn.pingCount == 3 {
									*cn.timeoutCount++
								} else {
									*cn.timeoutCount = 0
									*c.cluster[n].pingCount = 3
								}
							}
						}
					}
					c.lock.Unlock()
				}
			}
		}
	}()
	wch := cli.Watch(context.Background(), c.group, clientv3.WithPrefix())
	for wresp := range wch { //blocked
		for _, ev := range wresp.Events {
			cmds := strings.Split(string(ev.Kv.Key), "#")
			switch cmds[1] {
			case "ping":
				rnm := string(ev.Kv.Value)
				if rnm != c.local.Name {
					c.lock.Lock()
					v, exist := c.cluster[rnm]
					if exist {
						*v.pingCount--
					}
					c.lock.Unlock()
				}
			case "join":
				var rnd LocalNode
				err := json.Unmarshal(ev.Kv.Value, &rnd)
				if err == nil {
					cli.Put(context.Background(), c.group+"#joined", string(nd))
				}
			case "joined":
				var rnd LocalNode
				err := json.Unmarshal(ev.Kv.Value, &rnd)
				if err == nil {
					c.lock.Lock()
					_, joined := c.cluster[rnd.Name]
					if !joined {
						core.AppLog.Printf("Node [%s] has joined\n", rnd.Name)
						rnd.pingCount = new(int8)
						rnd.timeoutCount = new(uint8)
						*rnd.pingCount = 0
						*rnd.timeoutCount = 0
						c.cluster[rnd.Name] = rnd
						c.grouping()
						if c.started.Load() {
							c.clistener.MemberJoined(rnd.Node)
						}
					}
					c.lock.Unlock()
				}
			default:
				c.clistener.Updated(cmds[1], string(ev.Kv.Value))
			}
		}
	}
	core.AppLog.Printf("Cluster shut down [%s]\n", c.Group())
	return nil
}

func (c *EtcdCluster) Wait() {
	c.starting.Wait()
}

func (c *EtcdCluster) Started() {
	c.started.Store(true)
}
