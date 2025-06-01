package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"gameclustering.com/internal/util"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Node struct {
	Name         string `json:"name"`
	HttpEndpoint string `json:"http"`
	TcpEndpoint  string `json:"tcp"`
	pingCount    *int8  `json:"-"`
	timeoutCount *uint8 `json:"-"`
}

type Etc struct {
	Quit            chan bool
	Started         *sync.WaitGroup
	Group           string
	partitionNumber int
	EtcdEndpoints   []string
	local           Node
	lock            *sync.Mutex
	cluster         map[string]Node
	partition       []string
}

func NewEtc(group string, pNumber int, etcEndpoints []string, local Node) Etc {
	etc := Etc{Group: group, partitionNumber: pNumber, EtcdEndpoints: etcEndpoints, local: local}
	etc.lock = &sync.Mutex{}
	etc.cluster = make(map[string]Node)
	etc.partition = make([]string, pNumber)
	etc.Quit = make(chan bool)
	etc.Started = &sync.WaitGroup{}
	etc.Started.Add(1)
	return etc
}

func (c *Etc) Join() error {
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
			fmt.Printf("Waiting for member joining [%s]: %v\n", c.Group, w)
			if r == 0 {
				cli.Put(context.Background(), c.Group+"#join", string(nd))
			} else {
				cli.Put(context.Background(), c.Group+"#ping", string(nd))
			}
		}
		c.Started.Done()
	}()
	go func() {
		c.Started.Wait() //blocked
		tik := time.NewTicker(1 * time.Second)
		pct := 5
		for {
			select {
			case <-c.Quit:
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
								fmt.Printf("Node timeout %d %d %s %v\n", *cn.pingCount, *cn.timeoutCount, cn.Name, p)
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
			//fmt.Printf("Orignal %s\n", string(ev.Kv.Key))
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
				var rnd Node
				err := json.Unmarshal(ev.Kv.Value, &rnd)
				if err == nil {
					cli.Put(context.Background(), c.Group+"#joined", string(nd))
				}
			case "joined":
				var rnd Node
				err := json.Unmarshal(ev.Kv.Value, &rnd)
				if err == nil {
					c.lock.Lock()
					_, joined := c.cluster[rnd.Name]
					if !joined {
						fmt.Printf("Node [%s] has joined\n", rnd.Name)
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
				fmt.Printf("unwatch key %s\n", cmds[1])
			}
		}
	}
	fmt.Printf("Cluster shut down [%s]\n", c.Group)
	return nil
}
func (c *Etc) Local() Node {
	return c.local
}
func (c *Etc) View() iter.Seq[Node] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return maps.Values(c.cluster)
}

func (c *Etc) Partition(key []byte) Node {
	p := util.Partition(key, uint32(c.partitionNumber))
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cluster[c.partition[p]]
}

func (c *Etc) group() {
	sz := len(c.cluster)
	fmt.Printf("Cluster grouping %d\n", sz)
	nds := make([]string, sz)
	i := 0
	for n := range c.cluster {
		nds[i] = n
		i++
	}
	slices.Sort(nds)
	for p := range c.partitionNumber {
		i := p % sz
		c.partition[p] = nds[i]
		fmt.Printf("Partition %d %s %d\n", i, nds[i], p)
	}
}

func (c *Etc) Atomic(t Exec) error {
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
	mutex := concurrency.NewMutex(session, c.Group+"#lock")
	ctx := context.Background()
	mutex.Lock(ctx)
	defer mutex.Unlock(ctx)
	return t(&EtcdClient{cli: cli, prefix: c.Group})
}
