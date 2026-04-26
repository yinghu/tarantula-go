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
	revision uint64
	pending  []*protocol.Job
	jobIndex int
}


type TransactionResource struct {
	resource  *protocol.Transaction
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

	s *DataServiceProvider

	tasks   chan *protocol.Task
	jobs    chan *protocol.Job
	updates chan *protocol.Meta
}

func (m *TaskManager) set(t *TaskResource) {

	if t.resource.Validator != nil {
		t.pending = append(t.pending, t.resource.Validator)
	}
	t.pending = append(t.pending, t.resource.Job)
	t.jobIndex = 0
	job := t.pending[t.jobIndex]
	job.Meta.State = protocol.TCC_JOB_TIMEOUT
	m.schedule(t, job)
}

func (m *TaskManager) schedule(t *TaskResource, job *protocol.Job) {
	go m.s.updateTask(t, func() {
		core.AppLog.Debug().Msgf("task updated %d", t.revision)
	})
	go func() {
		m.jobs <- job
	}()
}

func (m *TaskManager) start(j *JobResource) {
	m.tms[j.resource.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(j.resource.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{TaskId: j.resource.Meta.TaskId, JobId: j.resource.Meta.Id, State: protocol.TCC_JOB_TIMEOUT}
	})}
	j.joinParties = len(j.resource.Transactions)
	j.confirmed = 0
	for _, tc := range j.resource.Transactions {
		j.joining[tc.Meta.Id] = &TransactionResource{resource: tc}
		tc.Meta.State = protocol.TCC_RESERVING
		tc.Meta.Time = timestamppb.Now()
		m.tms[tc.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: j.resource.Meta.TaskId, JobId: j.resource.Meta.Id, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		}), p: func() {
			core.AppLog.Debug().Msg("retry to reserve with timeout")
			tc.Meta.Time = timestamppb.Now()
			go m.s.runAskReserve(tc)
		}, d: time.Duration(tc.Meta.Timeout) * time.Second, r: tc.Meta.Retries}
		//ask to reserve
		go m.s.runAskReserve(tc)
	}
}

func (m *TaskManager) stop(t *JobResource) {
	tr := m.trs[t.resource.Meta.TaskId]
	tr.resource.Meta.State = protocol.TCC_FINISHED
	tr.resource.Validator.Meta.State = protocol.TCC_FINISHED
	tr.resource.Job.Meta.State = protocol.TCC_FINISHED
	m.end(tr)
}

func (m *TaskManager) confirmed(t *JobResource) {
	t.confirmed = 0
	for _, tc := range t.resource.Transactions {
		tc.Meta.State = protocol.TCC_CONFIRMED
		//retry to finish
		m.tms[tc.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: t.resource.Meta.TaskId, JobId: t.resource.Meta.JobId, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		}), p: func() {
			core.AppLog.Debug().Msg("retry to finish with confirm/timeout")
			go m.s.runAskFinish(tc.Meta)
		}, d: time.Duration(tc.Meta.Timeout) * time.Second, r: tc.Meta.Retries}

		//ask to finish
		go m.s.runAskFinish(tc.Meta)

	}
}

func (m *TaskManager) canceled(t *JobResource) {
	t.confirmed = 0
	tr := m.trs[t.resource.Meta.TaskId]
	tr.jobIndex = len(tr.pending)
	for _, tc := range t.resource.Transactions {
		m.closeTimer(tc.Meta.Id)
		tc.Meta.State = protocol.TCC_CANCELED
		//retry to finish on cancel
		m.tms[tc.Meta.Id] = &Timeout{t: time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: t.resource.Meta.TaskId, JobId: t.resource.Meta.Id, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		}), p: func() {
			core.AppLog.Debug().Msg("retry to finish with cancel/timeout")
			go m.s.runAskFinish(m.copy(tc.Meta))
		}, d: time.Duration(tc.Meta.Timeout) * time.Second, r: tc.Meta.Retries}
		//ask to finish
		go m.s.runAskFinish(m.copy(tc.Meta))
	}
}

func (m *TaskManager) finished(t *JobResource) {
	m.closeTimer(t.resource.Meta.Id)
	tr := m.trs[t.resource.Meta.TaskId]
	t.resource.Meta.State = protocol.TCC_FINISHED
	if tr.jobIndex+1 < len(tr.pending) {
		tr.jobIndex++
		next := tr.pending[tr.jobIndex]
		next.Meta.State = protocol.TCC_JOB_TIMEOUT
		m.schedule(tr, next)
		return
	}
	m.end(tr)
}

func (m *TaskManager) end(t *TaskResource) {
	t.resource.Meta.State = protocol.TCC_FINISHED

	go m.s.updateTask(t, func() {
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
	tr, err := m.s.load(meta.TaskId)
	if err != nil {
		core.AppLog.Warn().Msgf("task not existed %d", meta.TaskId)
		return nil, err
	}
	if tr.resource.Meta.State == protocol.TCC_FINISHED {
		return nil, fmt.Errorf("task alread finished")
	}
	m.trs[meta.TaskId] = tr
	if tr.resource.Validator.Meta.State != protocol.TCC_FINISHED {
		tr.pending = append(tr.pending, tr.resource.Validator)
	}
	if tr.resource.Job.Meta.State != protocol.TCC_FINISHED {
		tr.pending = append(tr.pending, tr.resource.Job)
	}
	if len(tr.pending) == 0 {
		return tr, fmt.Errorf("no job available")
	}
	tr.jobIndex = 0

	job := tr.pending[tr.jobIndex]
	job.Meta.State = protocol.TCC_JOB_TIMEOUT
	go m.s.updateTask(tr, func() {
		core.AppLog.Debug().Msgf("task updated from reload %d", tr.revision)
	})
	tj := JobResource{resource: job, joining: make(map[uint64]*TransactionResource)}
	m.tjs[job.Meta.Id] = &tj
	m.start(&tj)
	return tr, nil
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
	m.jobs = make(chan *protocol.Job, 10)
	for m.s.running {
		select {
		case task := <-m.tasks:
			tr := TaskResource{resource: task, revision: 1, pending: make([]*protocol.Job, 0)}
			m.trs[task.Meta.Id] = &tr
			m.set(&tr)
		case job := <-m.jobs:
			tj := JobResource{resource: job, joining: make(map[uint64]*TransactionResource)}
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
			tj := m.tjs[meta.JobId]
			switch meta.State {
			case protocol.TCC_CONFIRMED:
				m.closeTimer(meta.Id)
				if tj.join(meta) {
					m.confirmed(tj)
				}
			case protocol.TCC_CANCELED:
				m.closeTimer(meta.Id)
				m.canceled(tj)

			case protocol.TCC_FINISHED:
				m.closeTimer(meta.Id)
				if tj.join(meta) {
					m.finished(tj)
				}

			case protocol.TCC_TRANSACTION_TIMEOUT:
				core.AppLog.Debug().Msgf("task transaction timeout %d", tr.jobIndex)
				m.timeout(meta.Id, meta)
			case protocol.TCC_JOB_TIMEOUT:
				core.AppLog.Debug().Msgf("task job timeout %d", tr.jobIndex)
				m.timeout(meta.JobId, meta)
				m.stop(tj)
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
	cp := protocol.Meta{TaskId: meta.TaskId, JobId: meta.JobId, Id: meta.Id, State: meta.State, NodeId: meta.NodeId, Tag: meta.Tag, Name: meta.Name}
	cp.Timeout = meta.Timeout
	cp.Retries = meta.Retries
	cp.Time = meta.Time
	return &cp
}
