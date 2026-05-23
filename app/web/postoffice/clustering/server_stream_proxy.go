package clustering

import "google.golang.org/grpc"

type ServerStreamPorxy struct {
	grpc.ServerStream
}

func (r *ServerStreamPorxy) SendMsg(m any) error {
	return r.ServerStream.SendMsg(m)
}

func (r *ServerStreamPorxy) RecvMsg(m any) error {
	return r.ServerStream.RecvMsg(m)
}
