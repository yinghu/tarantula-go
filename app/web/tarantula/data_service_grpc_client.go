package main

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/tarantula/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (m *DataServiceProvider) ClientGet(target *core.Node, request *core.GetRequest) ([]byte, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	var dt *protocol.Data
	dt, err = dsp.Get(context.Background(), &protocol.Request{Key: request.Key})
	if err != nil {
		return nil, err
	}
	return dt.Value, nil
}

func (m *DataServiceProvider) ClientSet(target *core.Node, request *core.SetRequest) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	kv := protocol.Data{Key: request.Key, Value: request.Value}
	return dsp.Set(context.Background(), &kv)
}