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
				task.Meta.State = protocol.TCC_FINISHED
				m.updates <- task.Meta
			})}
			m.trs[task.Meta.TaskId] = &tr
			go tr.Start()
		case meta := <-m.updates:
			core.AppLog.Debug().Msgf("update %v", meta)
			//task := m.trs[meta.TaskId]
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
