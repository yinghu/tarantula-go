package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type LocalNode struct {
	core.Node
	pingCount    *int8  `json:"-"`
	timeoutCount *uint8 `json:"-"`
}

type EtcdCluster struct {
	kyl           core.KeyListener
	quit          chan bool
	started       *sync.WaitGroup
	Group         string
	EtcdEndpoints []string
	local         LocalNode
	lock          *sync.Mutex
	cluster       map[string]LocalNode
	partition     []string
}

func newCluster(group string, etcEndpoints []string, local LocalNode, kl core.KeyListener) core.Cluster {
	etc := EtcdCluster{Group: group, EtcdEndpoints: etcEndpoints, local: local}
	etc.kyl = kl
	etc.lock = &sync.Mutex{}
	etc.cluster = make(map[string]LocalNode)
	etc.partition = make([]string, core.CLUSTER_PARTITION_NUM)
	etc.quit = make(chan bool)
	etc.started = &sync.WaitGroup{}
	etc.started.Add(1)
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
			core.AppLog.Printf("Waiting for member joining [%s]: %v\n", c.Group, w)
			if r == 0 {
				cli.Put(context.Background(), c.Group+"#join", string(nd))
			} else {
				cli.Put(context.Background(), c.Group+"#ping", string(nd))
			}
		}
		c.started.Done()
	}()
	go func() {
		c.started.Wait() //blocked
		tik := time.NewTicker(1 * time.Second)
		pct := 5
		for {
			select {
			case <-c.quit:
				cli.Close()
				return
			case p := <-tik.C:
				cli.Put(context.Background(), c.Group+"#ping", c.local.Name)
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
								c.group()
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
	wch := cli.Watch(context.Background(), c.Group, clientv3.WithPrefix())
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
					cli.Put(context.Background(), c.Group+"#joined", string(nd))
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
						c.group()
					}
					c.lock.Unlock()
				}
			default:
				c.kyl.Updated(cmds[1], string(ev.Kv.Value))
			}
		}
	}
	core.AppLog.Printf("Cluster shut down [%s]\n", c.Group)
	return nil
}
func (c *EtcdCluster) Local() core.Node {
	return c.local.Node
}
func (c *EtcdCluster) View() []core.Node {
	c.lock.Lock()
	defer c.lock.Unlock()
	nv := make([]core.Node, 0)
	for _, v := range c.cluster {
		nv = append(nv, v.Node)
	}
	return nv
}

func (c *EtcdCluster) Partition(key []byte) core.Node {
	p := util.Partition(key, uint32(core.CLUSTER_PARTITION_NUM))
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cluster[c.partition[p]].Node
}

func (c *EtcdCluster) group() {
	sz := len(c.cluster)
	core.AppLog.Printf("Cluster grouping %d\n", sz)
	nds := make([]string, sz)
	i := 0
	for n := range c.cluster {
		nds[i] = n
		i++
	}
	slices.Sort(nds)
	for p := range core.CLUSTER_PARTITION_NUM {
		i := p % sz
		c.partition[p] = nds[i]
	}
}

func (c *EtcdCluster) Atomic(prefix string, t core.Exec) error {
	if prefix == "" {
		prefix = c.Group
		core.AppLog.Printf("Reset Lock prefix %s\n", prefix)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	session, err := concurrency.NewSession(cli)
	if err != nil {
		return err
	}
	defer session.Close()
	mutex := concurrency.NewMutex(session, prefix+"#lock")
	ctx := context.Background()
	mutex.Lock(ctx)
	defer mutex.Unlock(ctx)
	return t(&EtcdClient{cli: cli, prefix: prefix})
}

func (c *EtcdCluster) Wait() {
	c.started.Wait()
}

func (c *EtcdCluster) Quit() {
	c.quit <- true
}

func (c *EtcdCluster) OnJoin(join core.Node) {
	fmt.Printf("Node joined %v\n", join)
}

func (c *EtcdCluster) Listener() core.KeyListener {
	return c.kyl
}
