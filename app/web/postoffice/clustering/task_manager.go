package clustering

import (
	"fmt"
	"sync"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type TaskResource struct {
	*sync.Mutex
	resource *protocol.Task
	timer    *time.Timer
	ds       *DataServiceProvider
}

func (t *TaskResource) Update(u *protocol.Meta) {
	t.Lock()
	defer t.Unlock()
}

func (t *TaskResource) Start() {
	t.Lock()
	defer t.Unlock()
	for _, tc := range t.resource.Transactions {
		t.ds.runReserve(tc)
	}
}

func (t *TaskResource) Monitor() {
	t.Lock()
	defer t.Unlock()
	core.AppLog.Debug().Msgf("timeout %v", t.resource)
}

type TaskManager struct {
	*sync.Mutex
	trs map[uint64]*TaskResource
	s   *DataServiceProvider
}

func (m *TaskManager) Get(tid uint64) (*TaskResource, error) {
	m.Lock()
	defer m.Unlock()
	r, ok := m.trs[tid]
	if !ok {
		return nil, fmt.Errorf("not existed")
	}
	return r, nil
}

func (m *TaskManager) Set(t *protocol.Task) *TaskResource {
	m.Lock()
	defer m.Unlock()
	r, ok := m.trs[t.Meta.Id]
	if !ok {
		tr := TaskResource{resource: t, ds: m.s}
		m.trs[t.Meta.Id] = &tr
		r = &tr
		r.timer = time.AfterFunc(time.Duration(r.resource.Meta.Timeout)*time.Second, r.Monitor)
	}
	return r
}
