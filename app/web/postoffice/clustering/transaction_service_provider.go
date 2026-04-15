package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

const (
	TASK_TIMEOUT_SECONDS uint32 = 10
)

func (c *DataServiceProvider) Setup(ctx context.Context, task *protocol.Task) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("running setup on target node %v", task)
	return &protocol.Response{Successful: true, Message: ""}, nil
}

func (c *DataServiceProvider) Reserve(ctx context.Context, in *protocol.Transaction) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("running reserve on target node %v", in)
	c.DMessager <- &protocol.Mail{Transaction: in, Opt: core.TRANS_MAIL}
	return &protocol.Response{Successful: true, Message: "run task"}, nil
}

func (c *DataServiceProvider) Confirmed(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("running confirm on target node %v", in)
	in.State = protocol.TCC_CONFIRMED
	c.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: in}, Opt: core.TRANS_MAIL}
	//waiting all parties to comfirm or cancel
	return &protocol.Response{Successful: true, Message: "run task"}, nil
}

func (c *DataServiceProvider) Canceled(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("running cancel on target node %v", in)
	in.State = protocol.TCC_CANCELED
	c.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: in}, Opt: core.TRANS_MAIL}
	//waiting all parties to canceled
	return &protocol.Response{Successful: true, Message: "run task"}, nil
}

func (c *DataServiceProvider) Finished(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("running finish on target node %v", in)
	//c.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: in}, Opt: core.TRANS_MAIL}
	//waiting all parties to canceled
	return &protocol.Response{Successful: true, Message: "run task"}, nil
}
