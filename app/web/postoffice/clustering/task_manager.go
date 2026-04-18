package clustering

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type TaskResource struct {
	resource *protocol.Task
	timer    *time.Timer
	ds       *DataServiceProvider
}

func (t *TaskResource) Start() {
	for _, tc := range t.resource.Transactions {
		t.ds.runReserve(tc)
	}
}

type TaskManager struct {
	trs     map[uint64]*TaskResource
	s       *DataServiceProvider
	tasks   chan *protocol.Task
	updates chan *protocol.Meta
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
			tr := TaskResource{resource: task, ds: m.s, timer: time.AfterFunc(time.Duration(task.Meta.Timeout)*time.Second, func() {
				m.updates <- &protocol.Meta{TaskId: task.Meta.Id, State: protocol.TCC_FINISHED}
			})}
			m.trs[task.Meta.TaskId] = &tr
			go tr.Start()
		case meta := <-m.updates:
			core.AppLog.Debug().Msgf("update %v", meta)
			task, existing := m.trs[meta.TaskId]
			if !existing {
				task, err := m.s.load(meta.TaskId)
				if err != nil {
					core.AppLog.Warn().Msgf("task not existed %d", meta.TaskId)
					continue
				}
				m.trs[meta.TaskId] = &TaskResource{resource: task, timer: time.AfterFunc(time.Duration(task.Meta.Timeout)*time.Second, func() {
					m.updates <- &protocol.Meta{TaskId: task.Meta.Id, State: protocol.TCC_FINISHED}
				})}
			}
			core.AppLog.Debug().Msgf("task loaded %v", task)
			switch meta.State {
			case protocol.TCC_CONFIRMED:
				m.s.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: meta}, Opt: core.TRANS_MAIL}
			case protocol.TCC_CANCELED:
				m.s.DMessager <- &protocol.Mail{Transaction: &protocol.Transaction{Meta: meta}, Opt: core.TRANS_MAIL}
			case protocol.TCC_FINISHED:
				core.AppLog.Debug().Msgf("task finished %v", meta)
			}
		}
	}
}
