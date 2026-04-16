package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
)

const (
	TASK_TIMEOUT_SECONDS        uint32 = 30
	TRANSACTION_TIMEOUT_SECONDS uint32 = 10
	TASK_RETRY_MAX              uint32 = 3
)

func (c *DataServiceProvider) Setup(ctx context.Context, task *protocol.Task) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("running setup on target node %v", task)
	req := c.request(task.Meta.Id)
	get := GetData{Request: req}
	data, err := c.get(get)
	if err != nil {
		core.AppLog.Warn().Msgf("cannot load task %s", err.Error())
		return &protocol.Response{Successful: false}, err
	}
	tb := persistence.TaskBuilder{}
	tk, err := tb.From(data.Value)
	if err != nil {
		core.AppLog.Warn().Msgf("cannot parse task %s", err.Error())
		return &protocol.Response{Successful: false}, err
	}
	core.AppLog.Debug().Msgf("running setup on target node %v", tk)
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

func (c *DataServiceProvider) request(taskId uint64) *protocol.Request {
	buff := core.NewBuffer(8)
	buff.WriteUInt64(taskId)
	buff.Flip()
	k, _ := buff.Read(0)
	req := protocol.Request{Data: &protocol.Data{Key: k, Header: &protocol.Header{FactoryId: persistence.TASK_FACTORY_ID, ClassId: persistence.TASK_CLASS_ID}}}
	return &req
}
