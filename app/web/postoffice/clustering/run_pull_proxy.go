package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *DataServiceProvider) runPull(target string, set *protocol.Request) (grpc.ServerStreamingClient[protocol.Response], error) {
	core.AppLog.Debug().Msgf("run remote pull %s >= %d < %d", target, set.Prefix, set.Opt)
	tcp, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	return dsp.Pull(context.Background(), set)
}
