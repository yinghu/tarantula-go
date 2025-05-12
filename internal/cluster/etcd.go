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
	Quit          chan bool
	Started       *sync.WaitGroup
	Group         string
	EtcdEndpoints []string
	Local         Node
	lock          *sync.Mutex
	//partitions    map[uint32]string
	cluster map[string]Node
}

func NewEtc(group string, etcEndpoints []string, local Node) Etc {
	etc := Etc{Group: group, EtcdEndpoints: etcEndpoints, Local: local}
	etc.lock = &sync.Mutex{}
	//etc.partitions = make(map[uint32]string)
	etc.cluster = make(map[string]Node)
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
		tik := time.NewTicker(2 * time.Second)
		defer tik.Stop()
		for r := range 5 {
			w := <-tik.C
			fmt.Printf("Waiting for member joining : %v\n", w)
			if r == 0 {
				cli.Put(context.Background(), c.Group+"#join", string(nd))
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
							fmt.Printf("Check ping timeout %d %s %v\n", *cn.pingCount, cn.Name, p)
							*c.cluster[n].pingCount = 3
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
					//fmt.Println("Ping from [" + rnm + "][" + c.Local.Name + "]")
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
					fmt.Printf("Join from [%v]\n", rnd)
					cli.Put(context.Background(), c.Group+"#joined", string(nd))
				}
			case "joined":
				var rnd Node
				err := json.Unmarshal(ev.Kv.Value, &rnd)
				if err == nil {
					fmt.Printf("Joined from [%v]\n", rnd)
					c.lock.Lock()
					rnd.pingCount = new(int8)
					*rnd.pingCount = 3
					c.cluster[rnd.Name] = rnd
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
