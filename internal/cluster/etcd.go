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
}

type Etc struct {
	Quit          chan bool
	Started       *sync.WaitGroup
	Group         string
	EtcdEndpoints []string
	Local         Node
	lock          *sync.Mutex
	partitions    map[uint32]string
	cluster       map[string]Node
}

func NewEtc(group string, etcEndpoints []string, local Node) Etc {
	etc := Etc{Group: group, EtcdEndpoints: etcEndpoints, Local: local}
	etc.lock = &sync.Mutex{}
	etc.partitions = make(map[uint32]string)
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
			<-tik.C
			//fmt.Printf("Ticker : %d %v\n", c, t)
			if r == 0 {
				cli.Put(context.Background(), c.Group+"#join", string(nd))
			}
		}
		c.Started.Done()
	}()
	go func() {
		c.Started.Wait() //blocked
		tik := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-c.Quit:
				cli.Close()
				return
			case <-tik.C:
				//fmt.Printf("Ticker : %v\n", t)
				cli.Put(context.Background(), c.Group+"#ping", c.Local.Name)
			}
		}
	}()
	wch := cli.Watch(context.Background(), c.Group, clientv3.WithPrefix())
	for wresp := range wch { //blocked
		for _, ev := range wresp.Events {
			//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			cmds := strings.Split(string(ev.Kv.Key), "#")
			//fmt.Printf("Watching key : %s\n", cmds[1])
			switch cmds[1] {
			case "ping":
				rnm := string(ev.Kv.Value)
				if rnm != c.Local.Name {
					fmt.Println("Ping from [" + rnm + "][" + c.Local.Name + "]")
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
					c.cluster[rnd.Name] = rnd
					//JOIN
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
