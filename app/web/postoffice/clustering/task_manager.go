package clustering

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type TaskResource struct {
	resource          *protocol.Task
	timer             *time.Timer
	transactionTimers map[uint64]*time.Timer
}

type TaskManager struct {
	trs     map[uint64]*TaskResource
	s       *DataServiceProvider
	tasks   chan *protocol.Task
	updates chan *protocol.Meta
}

func (m *TaskManager) start(t *TaskResource) {
	t.timer = time.AfterFunc(time.Duration(t.resource.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{TaskId: t.resource.Meta.Id, State: protocol.TCC_TASK_TIMEOUT}
	})
	for _, tc := range t.resource.Transactions {
		t.transactionTimers[tc.Meta.Id] = time.AfterFunc(time.Duration(tc.Meta.Timeout)*time.Second, func() {
			m.updates <- &protocol.Meta{TaskId: t.resource.Meta.Id, Id: tc.Meta.Id, State: protocol.TCC_TRANSACTION_TIMEOUT}
		})
		go m.s.runReserve(tc)
	}
}

func (m *TaskManager) reload(meta *protocol.Meta) (*TaskResource, error) {
	task, err := m.s.load(meta.TaskId)
	if err != nil {
		core.AppLog.Warn().Msgf("task not existed %d", meta.TaskId)
		return nil, err
	}
	tr := TaskResource{resource: task, timer: time.AfterFunc(time.Duration(task.Meta.Timeout)*time.Second, func() {
		m.updates <- &protocol.Meta{TaskId: task.Meta.Id, State: protocol.TCC_TASK_TIMEOUT}
	})}
	m.trs[meta.TaskId] = &tr
	return &tr, nil
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
			m.trs[task.Meta.TaskId] = &tr
			go m.start(&tr)
		case meta := <-m.updates:
			core.AppLog.Debug().Msgf("update %v", meta)
			tr, existing := m.trs[meta.TaskId]
			if !existing {
				loaded, err := m.reload(meta)
				if err != nil {
					continue
				}
				tr = loaded
			}
			core.AppLog.Debug().Msgf("task loaded %v", tr)
			switch meta.State {
			case protocol.TCC_CONFIRMED:
				m.s.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: meta}, Opt: core.TRANS_MAIL}
			case protocol.TCC_CANCELED:
				m.s.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: meta}, Opt: core.TRANS_MAIL}
			case protocol.TCC_FINISHED:
				core.AppLog.Debug().Msgf("task finished %v", meta)
			case protocol.TCC_TRANSACTION_TIMEOUT:
				core.AppLog.Debug().Msgf("task transaction finished %v", meta)
			case protocol.TCC_TASK_TIMEOUT:
				core.AppLog.Debug().Msgf("task finished %v", meta)
			}
		}
	}
}
