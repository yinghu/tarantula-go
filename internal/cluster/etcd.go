package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Node struct {
	Name         string `json:"name"`
	HttpEndpoint string `json:"http"`
	TcpEndpoint  string `json:"tcp"`
	pingCount    *int8  `json:"-"`
}

type Etc struct {
	Quit            chan bool
	Started         *sync.WaitGroup
	Group           string
	PartitionNumber uint16
	EtcdEndpoints   []string
	Local           Node
	lock            *sync.Mutex
	cluster         map[string]Node
	partition       []string
}

func NewEtc(group string, partitionNumber uint16, etcEndpoints []string, local Node) Etc {
	etc := Etc{Group: group, PartitionNumber: partitionNumber, EtcdEndpoints: etcEndpoints, Local: local}
	etc.lock = &sync.Mutex{}
	etc.cluster = make(map[string]Node)
	etc.partition = make([]string, partitionNumber)
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
	nd, _ := json.Marshal(c.Local)
	go func() {
		tik := time.NewTicker(1 * time.Second)
		defer tik.Stop()
		for r := range 7 {
			w := <-tik.C
			fmt.Printf("Waiting for member joining : %v\n", w)
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
				cli.Put(context.Background(), c.Group+"#ping", c.Local.Name)
				pct--
				if pct == 0 {
					pct = 5
					c.lock.Lock()
					for n := range c.cluster {
						if n != c.Local.Name {
							cn := c.cluster[n]
							if *cn.pingCount == 3 {
								fmt.Printf("Node timeout %d %s %v\n", *cn.pingCount, cn.Name, p)
								delete(c.cluster, n)
								c.group()
							} else {
								fmt.Printf("RESET PING COUNT %d\n", *c.cluster[n].pingCount)
								*c.cluster[n].pingCount = 3
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
				if rnm != c.Local.Name {
					c.lock.Lock()
					v, exist := c.cluster[rnm]
					if exist {
						*v.pingCount--
						fmt.Printf("PING COUNT %d\n", *v.pingCount)
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
						*rnd.pingCount = 0
						c.cluster[rnd.Name] = rnd
						c.group()
					}
					c.lock.Unlock()
				}
			}
		}
	}
	fmt.Printf("Cluster shut down\n")
	return nil
}

func (c *Etc) View() iter.Seq[Node] {
	c.lock.Lock()
	defer c.lock.Unlock()
	return maps.Values(c.cluster)
}

func (c *Etc) group() {
	fmt.Printf("Cluster grouping %d\n", len(c.cluster))
}
