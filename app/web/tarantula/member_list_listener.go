package main

import (
	"fmt"
	"strings"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/hashicorp/memberlist"
)

const (
	REPLICA_RING_OPT int = 0

	ALL_RING_OPT int = 3

	ADD_NODE_OPT    = 5
	REMOVE_NODE_OPT = 6

	CLOSE_RING_OPT = 99

	NODE_STATE_LIVE     = 0
	NODE_STATE_DEAD     = 3
	NODE_STATE_SHUTDOWN = -1000
)

type RetryTrack struct {
	Err    error
	Reties int
	Suc    bool
}

type MemberListListener struct {
	MEvent    chan memberlist.NodeEvent
	MMerge    chan []core.Node
	MAlive    chan core.Node
	MPing     chan core.Node
	MConflict chan []core.Node
	MRequest  chan core.RingRequest
	*memberlist.Memberlist
	*MemberHashRing
	*DataServiceProvider
}

func (m *MemberListListener) toNode(e *memberlist.Node) core.Node {
	parts := strings.Split(e.Address(), ":")
	return core.Node{Name: e.Name, Meta: string(e.Meta), IP: fmt.Sprintf("%s:%d", parts[0], 7001), State: int(e.State)}
}

// event dispatch from event delegate
func (m *MemberListListener) Listen() {
	running := true
	for running {
		select {
		case e := <-m.MEvent:
			core.AppLog.Debug().Msgf("C EVET %v", e)
			switch e.Event {
			case memberlist.NodeJoin:
				m.OnAdd(m.toNode(e.Node))
			case memberlist.NodeLeave:
				if (m.LocalNode().Name) != e.Node.Name {
					m.OnRemove(m.toNode(e.Node))
				}
			case memberlist.NodeUpdate:
				m.OnUpdate(m.toNode(e.Node))
			}
		case mg := <-m.MMerge:
			m.OnMerge(mg)
		case ma := <-m.MAlive:
			m.OnLive(ma)
		case mp := <-m.MPing:
			m.OnPing(mp)
		case mc := <-m.MConflict:
			m.OnConflict(mc)
		case mr := <-m.MRequest:
			switch mr.Opt {
			case REPLICA_RING_OPT:
				nodes := m.keyRing(mr.Token, mr.Replicas)
				mr.Async <- nodes
			case ALL_RING_OPT:
				nodes := make([]core.Node, 0)
				for _, n := range m.nodes {
					nodes = append(nodes, n)
				}
				mr.Async <- nodes
			case ADD_NODE_OPT:
				mr.Async <- m.rangeNodeAdded(mr.Token)
			case REMOVE_NODE_OPT:
				mr.Async <- m.rangeNodeRemoved(mr.Token)
			case CLOSE_RING_OPT:
				running = false

			}
		}
	}
	core.AppLog.Info().Msg("local member listener has stopped")
}

func (m *MemberListListener) previousNode(r core.RingRequest) {
	m.MRequest <- r
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

func (m *MemberListListener) Get(get core.GetRequest) {
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

// delegate
func (m *MemberListListener) NodeMeta(limit int) []byte {
	//limit 512
	return nil //fmt.Appendf([]byte{}, "tarantula")
}

func (m *MemberListListener) NotifyMsg(msg []byte) {
	//broadcasting message
	core.AppLog.Printf("notify msf %s", string(msg))
}

func (m *MemberListListener) GetBroadcasts(overhead, limit int) [][]byte {
	//overhead 3 limit 1350
	return nil
}

func (m *MemberListListener) LocalState(join bool) []byte {
	if join {
		return []byte("mice")
	}
	return []byte("dog")
}
func (m *MemberListListener) MergeRemoteState(buf []byte, join bool) {
}

// ping delegate
func (m *MemberListListener) AckPayload() []byte {
	return []byte("cat")
}

func (m *MemberListListener) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	m.MPing <- m.toNode(other)
}

// merge delegate
func (m *MemberListListener) NotifyMerge(peers []*memberlist.Node) error {
	nodes := make([]core.Node, 0, len(peers))
	for _, n := range peers {
		nodes = append(nodes, m.toNode(n))
	}
	m.MMerge <- nodes
	return nil
}

// alive delegate
func (m *MemberListListener) NotifyAlive(peer *memberlist.Node) error {
	m.MAlive <- m.toNode(peer)
	return nil
}

// conflict delegate
func (m *MemberListListener) NotifyConflict(existing, other *memberlist.Node) {
	m.MConflict <- []core.Node{m.toNode(existing), m.toNode(other)}
}
