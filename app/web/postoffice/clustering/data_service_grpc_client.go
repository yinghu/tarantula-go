package clustering

import (
	context "context"
	"io"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Batch func(batch *protocol.DataSet)

func (m *DataServiceProvider) ClientGet(target *core.Node, request *core.DataRequest) ([]byte, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer tcp.Close()
	//dsp := protocol.NewDataServiceClient(tcp)
	var dt *protocol.Data
	//dt, err = dsp.Get(context.Background(), &protocol.Request{Key: request.Key, Header: &protocol.Header{FactoryId: request.FactoryId, ClassId: request.ClassId}})
	//if err != nil {
	//return nil, err
	//}
	return dt.Value, nil
}

func (m *DataServiceProvider) ClientSet(target *core.Node, request *core.DataRequest) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	//dsp := protocol.NewDataServiceClient(tcp)
	//kv := protocol.Data{Key: request.Key, Value: request.Value, Header: &protocol.Header{FactoryId: request.FactoryId, ClassId: request.ClassId}}
	//return dsp.Set(context.Background(), &kv)
	resp := protocol.Response{}
	return &resp, nil
}


func (m *DataServiceProvider) ClientUpdate(target *core.Node, request *core.DataRequest) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	//kv := protocol.Data{Key: request.Key, Value: request.Value}
	req := protocol.Request{}
	return dsp.Update(context.Background(), &req)
}

func (m *DataServiceProvider) ClientDelete(target *core.Node, request *core.DataRequest) (*protocol.Response, error) {
	tcp, err := grpc.NewClient(target.RpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	//kv := protocol.Data{Key: request.Key, Value: request.Value}
	req := protocol.Request{}
	return dsp.Delete(context.Background(), &req)
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
		resp := protocol.Response{Data: data.Data}
		batch(resp.Data)
	}
	return nil
}
