package clustering

import (
	"fmt"
	"sync"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type TaskResource struct {
	C        *sync.Mutex
	resource *protocol.Task
	timer    *time.Timer
	ds       *DataServiceProvider
}

func (t *TaskResource) Update(u *protocol.Meta) {
	t.C.Lock()
	defer t.C.Unlock()
}

func (t *TaskResource) Start() {
	t.C.Lock()
	defer t.C.Unlock()
	for _, tc := range t.resource.Transactions {
		t.ds.runReserve(tc)
	}
}

func (t *TaskResource) Monitor() {
	t.C.Lock()
	defer t.C.Unlock()
	core.AppLog.Debug().Msgf("timeout %v", t.resource)
}

type TaskManager struct {
	C   *sync.Mutex
	trs map[uint64]*TaskResource
	s   *DataServiceProvider
}

func (m *TaskManager) Get(tid uint64) (*TaskResource, error) {
	m.C.Lock()
	defer m.C.Unlock()
	r, ok := m.trs[tid]
	if !ok {
		return nil, fmt.Errorf("not existed")
	}
	return r, nil
}

func (m *TaskManager) Set(t *protocol.Task) *TaskResource {
	m.C.Lock()
	defer m.C.Unlock()
	r, ok := m.trs[t.Meta.Id]
	if !ok {
		tr := TaskResource{resource: t, ds: m.s, C: &sync.Mutex{}}
		m.trs[t.Meta.Id] = &tr
		r = &tr
		r.timer = time.AfterFunc(time.Duration(r.resource.Meta.Timeout)*time.Second, r.Monitor)
	}
	return r
}
