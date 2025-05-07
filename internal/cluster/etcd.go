package cluster

import (
	"context"
	"fmt"
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
	rch := cli.Watch(context.Background(), c.Group, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			cli.Close()
		}
	}
	return nil
}
