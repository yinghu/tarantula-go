package data

import (
	context "context"
	"net"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/grpc"
)

type DataServiceProvider struct {
	UnimplementedDataServiceServer
	Db *persistence.BadgerLocal
	Cs core.ClusterService
}

func (c *DataServiceProvider) Get(ctx context.Context, in *Request) (*Data, error) {
	data := Data{}
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

func (c *DataServiceProvider) Set(ctx context.Context, in *Data) (*Response, error) {
	err := c.Db.Db.Update(func(txn *badger.Txn) error {
		return txn.Set(in.Key, in.Value)
	})
	if err != nil {
		return &Response{Successful: false, Message: err.Error()}, err
	}
	return &Response{Successful: true, Message: "saved"}, nil
}

func (c *DataServiceProvider) Start() {
	c.Db.Open()
	tcp, err := net.Listen("tcp", ":7001")
	if err != nil {
		panic(err)
	}
	rpc := grpc.NewServer()
	RegisterDataServiceServer(rpc, c)
	core.AppLog.Printf("data service provider started on : %s", tcp.Addr().String())
	err = rpc.Serve(tcp)
	if err != nil {
		panic(err)
	}
}
