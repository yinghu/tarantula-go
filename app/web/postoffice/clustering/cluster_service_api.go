package clustering

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

func (m *MemberListListener) get(get core.DataRequest) {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}

	for retry.Reties > 0 {
		m.MRequest <- core.RingRequest{Token: m.RingToken(get.Key), Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		dt, err := m.DataServiceProvider.ClientGet(&ringNode, &get)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		get.Async <- core.Chunk{Remaining: false, Data: dt}
		retry.Suc = true
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	get.Async <- core.Chunk{Remaining: false, Data: []byte(retry.Err.Error())}
}
func (m *MemberListListener) Request(req core.DataRequest) {
	switch req.Opt {
	case core.GET_DATA_REQUEST:
		m.get(req)
	case core.UPDATE_DATA_REQUEST:
		m.update(req)
	case core.CREATE_DATA_REQUEST:
		//m.create(req)
	case core.DELETE_DATA_REQUEST:
		m.delete(req)
	case core.RESET_DATA_REQUEST:
		m.set(req)
	}
}

func (m *MemberListListener) update(set core.DataRequest) {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}

	for retry.Reties > 0 {
		m.MRequest <- core.RingRequest{Token: m.RingToken(set.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		_, err := m.DataServiceProvider.ClientSet(&ringNode, &set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		set.Async <- core.Chunk{Remaining: false, Data: util.ToJson(ringNode)}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			m.DataServiceProvider.ClientSet(&slave, &set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	set.Async <- core.Chunk{Remaining: false, Data: []byte(retry.Err.Error())}
}

func (m *MemberListListener) delete(set core.DataRequest) {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}

	for retry.Reties > 0 {
		m.MRequest <- core.RingRequest{Token: m.RingToken(set.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		_, err := m.DataServiceProvider.ClientSet(&ringNode, &set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		set.Async <- core.Chunk{Remaining: false, Data: util.ToJson(ringNode)}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			m.DataServiceProvider.ClientDelete(&slave, &set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	set.Async <- core.Chunk{Remaining: false, Data: []byte(retry.Err.Error())}
}

func (m *MemberListListener) set(set core.DataRequest) {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}

	for retry.Reties > 0 {
		m.MRequest <- core.RingRequest{Token: m.RingToken(set.Key), Replicas: REPLICA_MAX, Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		_, err := m.DataServiceProvider.ClientSet(&ringNode, &set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		set.Async <- core.Chunk{Remaining: false, Data: util.ToJson(ringNode)}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			m.DataServiceProvider.ClientSet(&slave, &set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	set.Async <- core.Chunk{Remaining: false, Data: []byte(retry.Err.Error())}
}
