package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Node struct {
	HttpEndpoint string `json:"http"`
	TcpEndpoint  string `json:"tcp"`
}

type Cluster struct {
	Group         string
	EtcdEndpoints []string
	Timeout       time.Duration
	Local         Node
}

func (c *Cluster) Join() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdEndpoints,
		DialTimeout: c.Timeout,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	tik := time.NewTicker(5 * time.Second)
	quit := make(chan bool)
	nd, _ := json.Marshal(c.Local)
	fmt.Printf("NODE : %q\n", nd)
	go func() {
		for {
			select {
			case <-quit:
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
			if cmds[1] == "quit" {
				cli.Close()
				quit <- true
			}
		}
	}
	tik.Stop()
	return nil
}
