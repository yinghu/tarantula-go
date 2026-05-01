package clustering

import (
	context "context"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
)

const (
	JOB_TIMEOUT_SECONDS         uint32 = 60
	TRANSACTION_TIMEOUT_SECONDS uint32 = 10
	TCC_RETRY_MAX               uint32 = 3
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

func (c *DataServiceProvider) load(taskId uint64) (*TaskResource, error) {
	buff := core.NewBuffer(8)
	buff.WriteUInt64(taskId)
	buff.Flip()
	k, _ := buff.Read(0)
	req := protocol.Request{Data: &protocol.Data{Key: k, Header: &protocol.Header{FactoryId: core.TASK_FACTORY_ID, ClassId: persistence.TASK_CLASS_ID}}}
	resp, err := c.runGet(&req)
	if err != nil {
		return nil, err
	}
	tb := persistence.TaskBuilder{}
	t, err := tb.From(resp.Data.List[0].Value)
	if err != nil {
		return nil, err
	}
	return &TaskResource{resource: t, revision: resp.Data.List[0].Header.Revision, pending: make([]*protocol.Job, 0)}, nil
}

func (c *DataServiceProvider) updateTask(t *TaskResource, clear ClearResource) {
	defer clear()
	tb := persistence.TaskBuilder{Target: t.resource}
	req, err := tb.Request()
	if err != nil {
		core.AppLog.Warn().Msgf("cannot request %s", err.Error())
		return
	}
	req.Data.Header.Revision = t.revision
	req.Opt = core.UPDATE_DATA_REQUEST
	resp, err := c.runUpdate(req)
	if err != nil {
		core.AppLog.Warn().Msgf("cannot update %s", err.Error())
		return
	}
	t.revision++
	core.AppLog.Info().Msgf("saved %v", resp)
	tx, err := c.load(t.resource.Meta.Id)
	if err != nil {
		core.AppLog.Error().Msgf("no task loaded %s", err.Error())
		return
	}
	core.AppLog.Debug().Msgf("v %v", tx.resource.Validator.Meta)
	core.AppLog.Debug().Msgf("j %v", tx.resource.Job.Meta)
	core.AppLog.Debug().Msgf("t %v", tx.resource.Meta)
}
