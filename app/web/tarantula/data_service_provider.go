package main

import (
	context "context"
	"net"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/tarantula/protocol"
	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataServiceProvider struct {
	protocol.UnimplementedDataServiceServer
	Db *persistence.BadgerLocal
	Cs core.ClusterService
	RNode <-chan []core.Node
}

func (c *DataServiceProvider) Get(ctx context.Context, in *protocol.Request) (*protocol.Data, error) {
	data := protocol.Data{}
	err := c.Db.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(in.Key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			data.Value = val
			return nil
		})
	})
	return &data, err
}

func (c *DataServiceProvider) Set(ctx context.Context, in *protocol.Data) (*protocol.Response, error) {
	err := c.Db.Db.Update(func(txn *badger.Txn) error {
		return txn.Set(in.Key, in.Value)
	})
	if err != nil {
		return &protocol.Response{Successful: false, Message: err.Error()}, err
	}
	return &protocol.Response{Successful: true, Message: "saved"}, nil
}

func (c *DataServiceProvider) Start() {
	c.Db.Open()
	tcp, err := net.Listen("tcp", ":7001")
	if err != nil {
		panic(err)
	}
	rpc := grpc.NewServer()
	protocol.RegisterDataServiceServer(rpc, c)
	core.AppLog.Printf("data service provider started on : %s", tcp.Addr().String())
	err = rpc.Serve(tcp)
	if err != nil {
		panic(err)
	}
}

func (m *DataServiceProvider) ClientGet(target *core.Node, request *core.GetRequest) ([]byte, error) {
	tcp, err := grpc.NewClient(target.IP, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	//core.AppLog.Printf("target node %s %s %d\n", ringNode.IP, ringNode.Name, ringNode.RingToken)
	tcp, err := grpc.NewClient(target.IP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	kv := protocol.Data{Key: request.Key, Value: request.Value}
	return dsp.Set(context.Background(), &kv)
}
func (m *DataServiceProvider) RingUpdated() {
	for nlist := range m.RNode {
		for _, n := range nlist {
			core.AppLog.Printf("node updated IP : %s NAME : %s RING TOKEN : %d STATE : %d", n.IP, n.Name, n.RingToken, n.State)
		}
	}
}
