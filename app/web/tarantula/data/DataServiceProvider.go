package data

import (
	context "context"
	"net"

	"gameclustering.com/internal/core"
	"google.golang.org/grpc"
)

type DataServiceProvider struct {
	UnimplementedDataServiceServer
}

func (c *DataServiceProvider) Get(ctx context.Context, in *Request) (*Data, error) {
	data := Data{Key: []byte("key"), Value: []byte("value")}
	core.AppLog.Println("calling from get")
	return &data, nil
}

func (c *DataServiceProvider) Set(ctx context.Context, in *Data) (*Response, error) {
	resp := Response{Successful: true, Code: 0, Message: "save"}
	core.AppLog.Println("calling from set")
	return &resp, nil
}

func (c *DataServiceProvider) Start() {
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
