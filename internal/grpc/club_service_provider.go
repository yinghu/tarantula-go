package grpc

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ClubServiceProvider struct {
	UnimplementedClubServiceServer
}

func (cs ClubServiceProvider) Club(context.Context, *emptypb.Empty) (*Club, error) {
	return nil, nil
}

func (cs ClubServiceProvider) Register(context.Context, *Person) (*Club, error) {
	return nil, nil
}

func (cs ClubServiceProvider) Unregister(context.Context, *Person) (*Club, error) {
	return nil, nil
}

func (cs ClubServiceProvider) List(*emptypb.Empty, grpc.ServerStreamingServer[Person]) error {
	return nil
}

func (cs ClubServiceProvider) Post(grpc.ClientStreamingServer[Person, Club]) error {
	return nil
}

func (cs ClubServiceProvider) Flow(grpc.BidiStreamingServer[Person, Club]) error {
	return nil
}
