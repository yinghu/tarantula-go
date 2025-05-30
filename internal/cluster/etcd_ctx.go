package cluster

import (
	"context"
	"errors"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	cli *clientv3.Client
}

func (c *EtcdClient) Put(key string, value string) error {
	ctx := context.Background()
	_, err := c.cli.Put(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}
func (c *EtcdClient) Get(key string) (string, error) {
	ctx := context.Background()
	r, err := c.cli.Get(ctx, key)
	if err != nil {
		return "", err
	}
	for _, ev := range r.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		return string(ev.Value),nil
	}
	return "",errors.New("no value getted")
}
func (c *EtcdClient) Del(key string) error {
	ctx := context.Background()
	r, err := c.cli.Delete(ctx, key)
	if err != nil {
		return err
	}
	if r.Deleted == 0{
		return errors.New("no value deleted")
	}
	return nil
}
