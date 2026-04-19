package clustering

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type TaskResource struct {
	resource *protocol.Task
	//transaction bookkeeping
	confirmed int
	finished  int
}

type TaskManager struct {
	trs     map[uint64]*TaskResource
	tms     map[uint64]*time.Timer
	s       *DataServiceProvider
	tasks   chan *protocol.Task
	trans   chan *protocol.Transaction
	updates chan *protocol.Meta
}

func (m *TaskManager) start(t *TaskResource) {
	m.tms[t.resource.Meta.Id] = time.AfterFunc(time.Duration(t.resource.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{TaskId: t.resource.Meta.Id, State: protocol.TCC_TASK_TIMEOUT}
	})
	t.confirmed = len(t.resource.Transactions)
	t.finished = t.confirmed
	for _, tc := range t.resource.Transactions {
		m.tms[tc.Meta.Id] = time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: t.resource.Meta.Id, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		})
		tc.Meta.State = protocol.TCC_RESERVING
		go m.s.runReserve(tc)
	}
}

func (m *TaskManager) confirmed(t *TaskResource) {
	for _, tc := range t.resource.Transactions {
		tc.Meta.State = protocol.TCC_CONFIRMED
		m.s.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: tc.Meta}, Opt: core.TRANS_MAIL}
	}
}

func (m *TaskManager) canceled(t *TaskResource) {
	for _, tc := range t.resource.Transactions {
		tc.Meta.State = protocol.TCC_CANCELED
		m.s.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: tc.Meta}, Opt: core.TRANS_MAIL}
	}
}

func (m *TaskManager) reload(meta *protocol.Meta) (*TaskResource, error) {
	core.AppLog.Debug().Msgf("reload task %d", meta.TaskId)
	task, err := m.s.load(meta.TaskId)
	if err != nil {
		core.AppLog.Warn().Msgf("task not existed %d", meta.TaskId)
		return nil, err
	}
	tr := TaskResource{resource: task}
	m.tms[meta.TaskId] = time.AfterFunc(time.Duration(task.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{TaskId: task.Meta.Id, State: protocol.TCC_TASK_TIMEOUT}
	})
	m.trs[meta.TaskId] = &tr
	return &tr, nil
}

func (m *TaskManager) Reserve(transaction *protocol.Transaction) {
	m.trans <- transaction
}

func (m *TaskManager) Update(meta *protocol.Meta) {
	m.updates <- meta
}

func (m *TaskManager) Set(t *protocol.Task) {
	m.tasks <- t
}

func (m *TaskManager) Wait() {
	m.tasks = make(chan *protocol.Task, 10)
	m.trans = make(chan *protocol.Transaction, 10)
	m.updates = make(chan *protocol.Meta, 10)
	for m.s.running {
		select {
		case task := <-m.tasks:
			tr := TaskResource{resource: task}
			m.trs[task.Meta.Id] = &tr
			m.start(&tr)
		case tran := <-m.trans:
			m.s.DMessager <- &protocol.Mail{Transaction: tran, Opt: core.TRANS_MAIL}
		case meta := <-m.updates:
			tr, existing := m.trs[meta.TaskId]
			if !existing {
				loaded, err := m.reload(meta)
				if err != nil {
					continue
				}
				tr = loaded
				core.AppLog.Debug().Msgf("task loaded %v", tr)
			}
			switch meta.State {
			case protocol.TCC_CONFIRMED:
				tr.confirmed--
				if tr.confirmed == 0 {
					m.confirmed(tr)
				}
				core.AppLog.Debug().Msgf("task confirmed %v", meta)
				
			case protocol.TCC_CANCELED:
				m.canceled(tr)
			case protocol.TCC_FINISHED:
				core.AppLog.Debug().Msgf("task finished %v", meta)
				tr.finished--
			case protocol.TCC_TRANSACTION_TIMEOUT:
				core.AppLog.Debug().Msgf("task transaction timeout %d %d", tr.confirmed, tr.finished)
			case protocol.TCC_TASK_TIMEOUT:
				core.AppLog.Debug().Msgf("task timeout %d %d", tr.confirmed, tr.finished)
			}
		}
	}
	clear(m.tms)
	clear(m.trs)
	close(m.tasks)
	close(m.trans)
	close(m.updates)
	core.AppLog.Warn().Msg("task manager stopped")
}
