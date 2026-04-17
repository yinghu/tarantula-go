package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
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
	tr := c.TManager.Set(task)
	tr.Start()
	return &protocol.Response{Successful: true, Meta: task.Meta}, nil
}

func (c *DataServiceProvider) Reserve(ctx context.Context, in *protocol.Transaction) (*protocol.Response, error) {
	core.AppLog.Debug().Msgf("00 running reserve on target node %v", in)
	tr, err := c.TManager.Get(in.Meta.TaskId)
	if err != nil {
		t, err := c.load(in.Meta.TaskId)
		if err != nil {
			return &protocol.Response{Successful: false}, err
		}
		tr = c.TManager.Set(t)
	}
	core.AppLog.Debug().Msgf("11 running reserve on target node %v", in)
	tr.Update(in.Meta)
	core.AppLog.Debug().Msgf("22 running reserve on target node %v", in)
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
	tf := event.NewTransactionEventFactory()
	e, _ := tf.FromTransactionEvent(&protocol.TransactionEvent{Meta: in})
	e.Event.Id = c.tid()
	c.runPublish(e)
	return &protocol.Response{Successful: true, Message: "run task"}, nil
}

func (c *DataServiceProvider) load(taskId uint64) (*protocol.Task, error) {
	buff := core.NewBuffer(8)
	buff.WriteUInt64(taskId)
	buff.Flip()
	k, _ := buff.Read(0)
	req := protocol.Request{Data: &protocol.Data{Key: k, Header: &protocol.Header{FactoryId: persistence.TASK_FACTORY_ID, ClassId: persistence.TASK_CLASS_ID}}}
	get := GetData{Request: &req}
	data, err := c.get(get)
	if err != nil {
		return nil, err
	}
	tb := persistence.TaskBuilder{}
	return tb.From(data.Value)
}
