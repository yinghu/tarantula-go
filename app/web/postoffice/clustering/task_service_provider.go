package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

func (c *DataServiceProvider) Run(ctx context.Context, in *protocol.Transaction) (*protocol.Response, error) {
	core.AppLog.Debug().Msg("running task on target node")
	c.DMessager <- &protocol.Mail{Transaction: in, Opt: core.TRANS_MAIL}
	return &protocol.Response{Successful: true, Message: "run task"}, nil
}
