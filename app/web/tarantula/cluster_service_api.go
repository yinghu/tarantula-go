package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

func (m *MemberListListener) Get(get core.GetRequest) {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}

	for retry.Reties > 0 {
		core.AppLog.Debug().Msgf("TRYING %d", retry.Reties)
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
		retry.Reties = 0
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	get.Async <- core.Chunk{Remaining: false, Data: []byte(retry.Err.Error())}
}

func (m *MemberListListener) Set(set core.SetRequest) {
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

func (m *MemberListListener) KeyRing(r core.RingRequest) {
	r.Replicas = REPLICA_MAX
	r.Opt = REPLICA_RING_OPT
	m.MRequest <- r
}

func (m *MemberListListener) HashRing(r core.RingRequest) {
	r.Opt = ALL_RING_OPT
	m.MRequest <- r
}
