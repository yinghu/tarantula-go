package cluster

import (
	"context"
	"encoding/json"
	"fmt"
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
	Started       sync.WaitGroup
	Group         string
	EtcdEndpoints []string
	Local         Node
}

func New(){
	
}

func (c *Etc) Join() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdEndpoints,
		DialTimeout: 5*time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	nd, _ := json.Marshal(c.Local)
	go func() {
		tik := time.NewTicker(2 * time.Second)
		defer tik.Stop()
		for c := range 5 {
			t := <-tik.C
			fmt.Printf("Ticker : %d %v\n", c, t)
			if c == 0 {
				cli.Put(context.Background(), "tarantula#join", string(nd))
			}
		}
		c.Started.Done()
	}()
	go func() {
		c.Started.Wait() //blocked
		tik := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-c.Quit:
				cli.Close()
				return
			case t := <-tik.C:
				fmt.Printf("Ticker : %v\n", t)
				cli.Put(context.Background(), "tarantula#ping", string(nd))
			}
		}
	}()
	wch := cli.Watch(context.Background(), c.Group, clientv3.WithPrefix())
	for wresp := range wch { //blocked
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			cmds := strings.Split(string(ev.Kv.Key), "#")
			fmt.Printf("Watching key : %s\n", cmds[1])
		}
	}
	fmt.Printf("Cluster shut down\n")
	return nil
}
