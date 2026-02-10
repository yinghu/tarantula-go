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
	Cs core.ClusterService
}

func (c *DataServiceProvider) Get(ctx context.Context, in *Request) (*Data, error) {
	ret, err := c.Db.Get(in.Database, in.Key)
	if err != nil {
		return nil, err
	}
	data := Data{Value: ret}
	id, _ := c.Cs.Id()
	core.AppLog.Printf("calling from get %d\n", id)
	return &data, nil
}

func (c *DataServiceProvider) Set(ctx context.Context, in *Data) (*Response, error) {
	err := c.Db.Put(in.Database, in.Key, in.Value)
	if err != nil {
		return &Response{Successful: false, Code: 100, Message: err.Error()}, err
	}
	id, _ := c.Cs.Id()
	core.AppLog.Printf("calling from set %d\n", id)
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
	core.AppLog.Printf("data service provider started on : %s\n", tcp.Addr().String())
	err = rpc.Serve(tcp)
	if err != nil {
		panic(err)
	}
}
