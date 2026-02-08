package data

import (
	context "context"
	"net"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"google.golang.org/grpc"
)

type DataServiceProvider struct {
	UnimplementedDataServiceServer
	Db *persistence.LMDBLocal
}

func (c *DataServiceProvider) Get(ctx context.Context, in *Request) (*Data, error) {
	ret, err := c.Db.Get(in.Database, in.Key)
	if err != nil {
		return nil, err
	}
	data := Data{Value: ret}
	core.AppLog.Println("calling from get")
	return &data, nil
}

func (c *DataServiceProvider) Set(ctx context.Context, in *Data) (*Response, error) {
	err := c.Db.Put(in.Database, in.Key, in.Value)
	if err != nil {
		return &Response{Successful: false, Code: 100, Message: err.Error()}, err
	}
	core.AppLog.Println("calling from set")
	return &Response{Successful: true}, nil
}

func (c *DataServiceProvider) Start() {
	c.Db.Open()
	tcp, err := net.Listen("tcp", ":7001")
	if err != nil {
		panic(err)
	}
	rpc := grpc.NewServer()
	RegisterDataServiceServer(rpc, c)
	core.AppLog.Printf("data service provider started on : %s\n", tcp.Addr().String())
	err = rpc.Serve(tcp)
	if err != nil {
		panic(err)
	}
}
