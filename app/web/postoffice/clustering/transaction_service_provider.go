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
	c.TManager.Set(task)
	return &protocol.Response{Successful: true, Meta: task.Meta}, nil
}

func (c *DataServiceProvider) Reserve(ctx context.Context, in *protocol.Transaction) (*protocol.Response, error) {
	c.TManager.Update(in.Meta)
	c.DMessager <- &protocol.Mail{Transaction: in, Opt: core.TRANS_MAIL}
	return &protocol.Response{Successful: true, Message: "run reserve task"}, nil
}

func (c *DataServiceProvider) Confirmed(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	in.State = protocol.TCC_CONFIRMED
	c.TManager.Update(in)

	return &protocol.Response{Successful: true, Message: "run confirm task"}, nil
}

func (c *DataServiceProvider) Canceled(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	in.State = protocol.TCC_CANCELED
	c.TManager.Update(in)
	return &protocol.Response{Successful: true, Message: "run cancel task"}, nil
}

func (c *DataServiceProvider) Finished(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	in.State = protocol.TCC_FINISHED
	c.TManager.Update(in)
	tf := event.NewTransactionEventFactory()
	e, _ := tf.FromTransactionEvent(&protocol.TransactionEvent{Meta: in})
	e.Event.Id = c.tid()
	c.runPublish(e)
	return &protocol.Response{Successful: true, Message: "run finish task"}, nil
}

func (c *DataServiceProvider) Xload(taskId uint64) (*protocol.Task, error) {
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
