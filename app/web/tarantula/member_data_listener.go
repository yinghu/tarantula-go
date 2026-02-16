package main

import (
	"context"

	"gameclustering.com/internal/core"
	"gameclustering.com/tarantula/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MemberDataListener struct {
	RNode <-chan []core.Node
}

func (m *MemberDataListener) RingUpdated() {
	for nlist := range m.RNode {
		for _, n := range nlist {
			core.AppLog.Printf("node updated IP : %s NAME : %s RING TOKEN : %d STATE : %d", n.IP, n.Name, n.RingToken, n.State)
		}
	}
}

func (m *MemberDataListener) Get(target *core.Node, request *core.GetRequest) ([]byte, error) {
	tcp, err := grpc.NewClient(target.IP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer tcp.Close()
	dsp := data.NewDataServiceClient(tcp)
	var dt *data.Data
	dt, err = dsp.Get(context.Background(), &data.Request{Database: request.Database, Key: request.Key})
	if err!=nil{
		return nil,err
	}
	return dt.Value, nil
}

func (m *MemberDataListener) Set(target *core.Node, request *core.SetRequest) (*data.Response, error) {
	//core.AppLog.Printf("target node %s %s %d\n", ringNode.IP, ringNode.Name, ringNode.RingToken)
	tcp, err := grpc.NewClient(target.IP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &data.Response{}, err
	}
	defer tcp.Close()
	dsp := data.NewDataServiceClient(tcp)
	kv := data.Data{Database: request.Database, Key: request.Key, Value: request.Value}
	return dsp.Set(context.Background(), &kv)
}
