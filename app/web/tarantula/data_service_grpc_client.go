package main

import (
	context "context"
	"io"

	"gameclustering.com/internal/core"
	"gameclustering.com/tarantula/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Batch func(batch *protocol.DataBatch)

func (m *DataServiceProvider) ClientGet(target *core.Node, request *core.GetRequest) ([]byte, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	var dt *protocol.Data
	dt, err = dsp.Get(context.Background(), &protocol.Request{Key: request.Key,Header: &protocol.Header{FactoryId: request.FactoryId,ClassId: request.ClassId}})
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
	kv := protocol.Data{Key: request.Key, Value: request.Value,Header: &protocol.Header{FactoryId: request.FactoryId,ClassId:request.ClassId}}
	return dsp.Set(context.Background(), &kv)
}

func (m *DataServiceProvider) ClientCreate(target *core.Node, request *core.SetRequest) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	kv := protocol.Data{Key: request.Key, Value: request.Value}
	return dsp.Create(context.Background(), &kv)
}

func (m *DataServiceProvider) ClientUpdate(target *core.Node, request *core.SetRequest) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	kv := protocol.Data{Key: request.Key, Value: request.Value}
	return dsp.Update(context.Background(), &kv)
}

func (m *DataServiceProvider) ClientDelete(target *core.Node, request *core.SetRequest) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	kv := protocol.Data{Key: request.Key, Value: request.Value}
	return dsp.Delete(context.Background(), &kv)
}

func (m *DataServiceProvider) ClientPull(target string, hash uint32, batch Batch) error {
	tcp, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	stream, err := dsp.Pull(context.Background(), &protocol.Request{Prefix: hash})
	if err != nil {
		return err
	}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("streaming error %s", err.Error())
			return err
		}
		batch(data)
	}
	return nil
}
