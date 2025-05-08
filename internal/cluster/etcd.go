package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Cluster struct {
	Group     string
	Endpoints []string
	Timeout   time.Duration
}

func (c *Cluster) Watch() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.Endpoints,
		DialTimeout: c.Timeout,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	wch := cli.Watch(context.Background(), c.Group, clientv3.WithPrefix())
	for wresp := range wch {//blocked
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			cmds := strings.Split(string(ev.Kv.Key), "#")
			if cmds[1] == "quit" {
				cli.Close()
			}
		}
	}
	return nil
}
