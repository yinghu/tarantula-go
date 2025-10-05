package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

func (c *EtcdClient) AppIndex(env string) AppIndex {
	apps := AppIndex{Index: make([]int,0), Env: env}
	ctx := context.Background()
	k := fmt.Sprintf("%s.%s", c.Prefix, env)
	data, err := c.Cli.Get(ctx, k)
	if err != nil {
		return apps
	}
	for _, ev := range data.Kvs {
		err = json.Unmarshal(ev.Value, &apps)
		if err != nil {
			fmt.Printf("app index parse error %s\n", err.Error())
		}
	}
	return apps
}

func (c *EtcdClient) SaveAppIndex(apps AppIndex) error{
	data,err :=json.Marshal(apps)
	if err!=nil{
		return err
	}
	ctx := context.Background()
	k := fmt.Sprintf("%s.%s", c.Prefix, apps.Env)
	_ ,err = c.Cli.Put(ctx,k,string(data))
	if err!=nil{
		return err
	}
	return nil
}

