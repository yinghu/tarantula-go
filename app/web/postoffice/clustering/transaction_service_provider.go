package clustering

import (
	context "context"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
)

const (
	TASK_TIMEOUT_SECONDS        uint32 = 60
	TRANSACTION_TIMEOUT_SECONDS uint32 = 10
	TASK_RETRY_MAX              uint32 = 3
)

type ClearResource func()

func (c *DataServiceProvider) Setup(ctx context.Context, task *protocol.Task) (*protocol.Response, error) {
	c.TManager.Set(task)
	return &protocol.Response{Successful: true, Meta: task.Meta}, nil
}

func (c *DataServiceProvider) AskReserve(ctx context.Context, in *protocol.Transaction) (*protocol.Response, error) {
	in.Meta.State = protocol.TCC_RESERVING
	c.TManager.Reserve(in)
	return &protocol.Response{Successful: true, Meta: &protocol.Meta{Name: in.Meta.Name}}, nil
}

func (c *DataServiceProvider) AskFinish(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	c.TManager.Finish(in)
	return &protocol.Response{Successful: true, Meta: &protocol.Meta{Name: in.Name}}, nil
}

func (c *DataServiceProvider) Confirmed(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	in.State = protocol.TCC_CONFIRMED
	c.TManager.Update(in)

	return &protocol.Response{Successful: true, Meta: &protocol.Meta{Name: in.Name}}, nil
}

func (c *DataServiceProvider) Canceled(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	in.State = protocol.TCC_CANCELED
	c.TManager.Update(in)
	return &protocol.Response{Successful: true, Meta: &protocol.Meta{Name: in.Name}}, nil
}

func (c *DataServiceProvider) Finished(ctx context.Context, in *protocol.Meta) (*protocol.Response, error) {
	in.State = protocol.TCC_FINISHED
	c.TManager.Update(in)
	return &protocol.Response{Successful: true, Meta: &protocol.Meta{Name: in.Name}}, nil
}

func (c *DataServiceProvider) load(taskId uint64) (*protocol.Task, error) {
	buff := core.NewBuffer(8)
	buff.WriteUInt64(taskId)
	buff.Flip()
	k, _ := buff.Read(0)
	req := protocol.Request{Data: &protocol.Data{Key: k, Header: &protocol.Header{FactoryId: persistence.TASK_FACTORY_ID, ClassId: persistence.TASK_CLASS_ID}}}
	resp, err := c.runGet(&req)
	if err != nil {
		return nil, err
	}
	tb := persistence.TaskBuilder{}
	return tb.From(resp.Data.List[0].Value)
}
func (c *DataServiceProvider) saveLog(meta *protocol.Meta) {
	tf := event.NewTransactionEventFactory()
	e, _ := tf.FromTransactionEvent(&protocol.TransactionEvent{Meta: meta})
	buff := core.NewBuffer(20)
	buff.WriteUInt64(meta.Id)
	buff.WriteUInt32(meta.State)
	buff.Flip()
	k, _ := buff.Read(0)
	e.Event.Key.Array = k
	go c.runPublish(e)
	go c.loadLog(meta)
}

func (c *DataServiceProvider) loadLog(meta *protocol.Meta) {
	time.Sleep(1 * time.Second)
	tf := event.NewTransactionEventFactory()
	buff := core.NewBuffer(20)
	buff.WriteUInt64(meta.Id)
	buff.WriteUInt32(meta.State)
	buff.Flip()
	k, _ := buff.Read(0)
	req := tf.GetRequest(k)
	resp, err := c.runGet(req)
	if err == nil && resp.Successful {
		topic, err := tf.Topic(resp.Data.List[0].Value)
		if err == nil {
			obj, _ := tf.Message(topic)
			te, ok := obj.(*protocol.TransactionEvent)
			if ok {
				core.AppLog.Debug().Msgf("META : %v", te.Meta)
			}
		}
	}
	core.AppLog.Debug().Msgf("LOG :%v", resp)
}

func (c *DataServiceProvider) updateTask(t *protocol.Task, clear ClearResource) {
	defer clear()
	tb := persistence.TaskBuilder{Target: t}
	req, err := tb.Request()
	if err != nil {
		core.AppLog.Warn().Msgf("cannot request %s", err.Error())
		return
	}
	req.Data.Header.Revision = 1
	req.Opt = core.UPDATE_DATA_REQUEST
	resp, err := c.runUpdate(req)
	if err != nil {
		core.AppLog.Warn().Msgf("cannot update %s", err.Error())
		return
	}
	core.AppLog.Info().Msgf("saved %v", resp)
}
