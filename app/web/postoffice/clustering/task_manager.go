package clustering

import (
	"fmt"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskResource struct {
	resource *protocol.Task
	//transaction bookkeeping
	confirmed int
	finished  int
	jobIndex  int
}

type JobResource struct {
	resource  *protocol.Job
	confirmed int
	finished  int
}

type Retrying func()
type Timeout struct {
	t *time.Timer
	d time.Duration
	r uint32
	p Retrying
}

type TaskManager struct {
	trs map[uint64]*TaskResource
	tjs map[uint64]*JobResource
	tms map[uint64]*Timeout

	s       *DataServiceProvider
	tasks   chan *protocol.Task
	jobs    chan *protocol.Job
	updates chan *protocol.Meta
}

func (m *TaskManager) set(t *TaskResource) {
	m.tms[t.resource.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(t.resource.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{TaskId: t.resource.Meta.Id, State: protocol.TCC_TASK_TIMEOUT}
	})}
	t.confirmed = len(t.resource.Jobs)
	t.finished = t.confirmed
	t.jobIndex = 0
	job := t.resource.Jobs[t.jobIndex]
	go func() {
		m.jobs <- job
	}()
}

func (m *TaskManager) start(j *JobResource) {
	m.tms[j.resource.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(j.resource.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{JobId: j.resource.Meta.Id, State: protocol.TCC_JOB_TIMEOUT}
	})}
	for _, tc := range j.resource.Transactions {
		tc.Meta.State = protocol.TCC_RESERVING
		tc.Meta.Time = timestamppb.Now()
		m.tms[tc.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: j.resource.Meta.Id, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		}), p: func() {
			core.AppLog.Debug().Msg("retry to reserve with timeout")
			tc.Meta.Time = timestamppb.Now()
			go m.s.runAskReserve(tc)
		}, d: time.Duration(tc.Meta.Timeout) * time.Second, r: tc.Meta.Retries}
		//ask to reserve
		go m.s.runAskReserve(tc)
	}
}

func (m *TaskManager) confirmed(t *TaskResource) {
	for _, tc := range t.resource.Jobs[0].Transactions {
		tc.Meta.State = protocol.TCC_CONFIRMED
		//retry to finish
		m.tms[tc.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: t.resource.Meta.Id, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		}), p: func() {
			core.AppLog.Debug().Msg("retry to finish with confirm/timeout")
			go m.s.runAskFinish(tc.Meta)
		}, d: time.Duration(tc.Meta.Timeout) * time.Second, r: tc.Meta.Retries}

		//ask to finish
		go m.s.runAskFinish(tc.Meta)

	}
}

func (m *TaskManager) canceled(t *TaskResource) {
	for _, tc := range t.resource.Jobs[0].Transactions {
		m.closeTimer(tc.Meta.Id)
		tc.Meta.State = protocol.TCC_CANCELED
		//retry to finish on cancel
		m.tms[tc.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: t.resource.Meta.Id, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		}), p: func() {
			core.AppLog.Debug().Msg("retry to finish with cancel/timeout")
			go m.s.runAskFinish(m.copy(tc.Meta))
		}, d: time.Duration(tc.Meta.Timeout) * time.Second, r: tc.Meta.Retries}
		//ask to finish
		go m.s.runAskFinish(m.copy(tc.Meta))
	}
}

func (m *TaskManager) finished(t *TaskResource) {
	t.resource.Meta.State = protocol.TCC_FINISHED
	m.closeTimer(t.resource.Meta.Id)
	go m.s.updateTask(t.resource, func() {
		m.updates <- &protocol.Meta{Id: t.resource.Meta.Id, State: protocol.TCC_TASK_CLEAR}
	})
	tf := event.NewTransactionEventFactory()
	e, _ := tf.FromTransactionEvent(&protocol.TransactionEvent{Meta: t.resource.Meta, Start: t.resource.Meta.Time, End: timestamppb.Now()})
	e.Event.Key.Array = core.ToBytes(m.s.seq)
	go m.s.runPublish(e)
}

func (m *TaskManager) closeTimer(mkey uint64) {
	tm, ok := m.tms[mkey]
	if !ok {
		return
	}
	tm.t.Stop()
	delete(m.tms, mkey)
}

func (m *TaskManager) timeout(mkey uint64, meta *protocol.Meta) {
	tm, ok := m.tms[mkey]
	if !ok {
		return
	}
	if tm.d > 0 && tm.r > 0 {
		core.AppLog.Debug().Msgf("retried %d", tm.r)
		// retry
		tm.t = time.AfterFunc(tm.d, func() {
			m.updates <- meta
		})
		tm.r--
		tm.p()
		return
	}
	delete(m.tms, mkey)

}

func (m *TaskManager) clearResource(rkey uint64) {
	delete(m.trs, rkey)
	core.AppLog.Debug().Msgf("task removed %d %d %d", rkey, len(m.tms), len(m.trs))
}

func (m *TaskManager) reload(meta *protocol.Meta) (*TaskResource, error) {
	core.AppLog.Debug().Msgf("reload task %d", meta.TaskId)
	task, err := m.s.load(meta.TaskId)
	if err != nil {
		core.AppLog.Warn().Msgf("task not existed %d", meta.TaskId)
		return nil, err
	}
	if task.Meta.State == protocol.TCC_FINISHED {
		return nil, fmt.Errorf("task alread finished")
	}
	tr := TaskResource{resource: task}
	m.tms[meta.TaskId] = &Timeout{t: time.AfterFunc(time.Duration(task.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{TaskId: task.Meta.Id, State: protocol.TCC_TASK_TIMEOUT}
	}), p: func() {
		core.AppLog.Debug().Msg("running task with timeout")
	}}
	m.trs[meta.TaskId] = &tr
	return &tr, nil
}

func (m *TaskManager) log(meta *protocol.Meta) {
	go m.s.saveLog(meta)
}

func (m *TaskManager) Reserve(transaction *protocol.Transaction) {
	m.s.DMessager <- &protocol.Mail{Transaction: transaction, Opt: core.TRANS_MAIL}
}
func (m *TaskManager) Finish(meta *protocol.Meta) {
	m.s.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: meta}, Opt: core.TRANS_MAIL}
}

func (m *TaskManager) Update(meta *protocol.Meta) {
	m.updates <- meta
}

func (m *TaskManager) Set(t *protocol.Task) {
	m.tasks <- t
}

func (m *TaskManager) Wait() {
	m.tasks = make(chan *protocol.Task, 10)
	m.updates = make(chan *protocol.Meta, 10)

	for m.s.running {
		select {
		case task := <-m.tasks:
			tr := TaskResource{resource: task}
			m.trs[task.Meta.Id] = &tr
			m.set(&tr)
		case job := <-m.jobs:
			tj := JobResource{resource: job}
			m.tjs[job.Meta.Id] = &tj
			m.start(&tj)
		case meta := <-m.updates:
			if meta.State == protocol.TCC_TASK_CLEAR {
				m.clearResource(meta.Id)
				continue
			}
			tr, existing := m.trs[meta.TaskId]
			if !existing {
				loaded, err := m.reload(meta)
				if err != nil {
					core.AppLog.Warn().Msgf("task load error %s", err.Error())
					continue
				}
				tr = loaded
				core.AppLog.Debug().Msgf("task loaded %v", tr.resource)
			}
			if tr.resource.Meta.State == protocol.TCC_FINISHED {
				core.AppLog.Warn().Msgf("task already finished %v", meta)
				continue
			}
			m.log(m.copy(meta))
			switch meta.State {
			case protocol.TCC_CONFIRMED:
				tr.confirmed--
				m.closeTimer(meta.Id)
				if tr.confirmed == 0 {
					m.confirmed(tr)
				}
				core.AppLog.Debug().Msgf("task confirmed %v", meta)
			case protocol.TCC_CANCELED:
				core.AppLog.Debug().Msgf("task canceled %v", meta)
				m.canceled(tr)
			case protocol.TCC_FINISHED:
				core.AppLog.Debug().Msgf("task finished %v", meta)
				tr.finished--
				m.closeTimer(meta.Id)
				if tr.finished == 0 {
					m.finished(tr)
				}

			case protocol.TCC_TRANSACTION_TIMEOUT:
				core.AppLog.Debug().Msgf("task transaction timeout %d %d", tr.confirmed, tr.finished)
				m.timeout(meta.Id, meta)
			case protocol.TCC_JOB_TIMEOUT:
				core.AppLog.Debug().Msgf("task job timeout %d %d", tr.confirmed, tr.finished)
				m.timeout(meta.JobId, meta)
			case protocol.TCC_TASK_TIMEOUT:
				core.AppLog.Debug().Msgf("task timeout %d %d", tr.confirmed, tr.finished)
				m.timeout(meta.TaskId, meta)
				m.finished(tr) //forcefully finished
			}
		}
	}
	clear(m.tms)
	clear(m.trs)
	close(m.tasks)
	close(m.updates)
	core.AppLog.Warn().Msg("task manager stopped")
}

func (m *TaskManager) copy(meta *protocol.Meta) *protocol.Meta {
	cp := protocol.Meta{TaskId: meta.TaskId, Id: meta.Id, State: meta.State, NodeId: meta.NodeId, Tag: meta.Tag, Name: meta.Name}
	cp.Timeout = meta.Timeout
	cp.Retries = meta.Retries
	cp.Time = meta.Time
	return &cp
}
