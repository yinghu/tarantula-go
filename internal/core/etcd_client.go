package core

import (
	"context"
	"errors"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	Cli    *clientv3.Client
	Prefix string
}

func (c *EtcdClient) Put(key string, value string) error {
	ctx := context.Background()
	_, err := c.Cli.Put(ctx, c.Prefix+"#"+key, value)
	if err != nil {
		return err
	}
	return nil
}
func (c *EtcdClient) Get(key string) (string, error) {
	ctx := context.Background()
	r, err := c.Cli.Get(ctx, c.Prefix+"#"+key)
	if err != nil {
		return "", err
	}
	for _, ev := range r.Kvs {
		return string(ev.Value), nil
	}
	return "", errors.New("no value getted")
}
func (c *EtcdClient) Del(key string) error {
	ctx := context.Background()
	r, err := c.Cli.Delete(ctx, c.Prefix+"#"+key)
	if err != nil {
		return err
	}
	if r.Deleted == 0 {
		return errors.New("no value deleted")
	}
	return nil
}

func (c *EtcdClient) List(prefix string, loaded KVLoad) error {

	ctx := context.Background()
	r, err := c.Cli.Get(ctx, c.Prefix+"#"+prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, ev := range r.Kvs {
		keys := strings.Split(string(ev.Key), "#")
		loaded(keys[1], string(ev.Value))
	}
	return nil
}