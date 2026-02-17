package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
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
	*MemberDataListener
	ct int
}

func (m *MemberListListener) toNode(e *memberlist.Node) core.Node {
	parts := strings.Split(e.Address(), ":")
	return core.Node{Name: e.Name, Meta: string(e.Meta), IP: fmt.Sprintf("%s:%d", parts[0], 7001), State: int(e.State)}
}

// event dispatch from event delegate
func (m *MemberListListener) Listen() {
	for {
		select {
		case e := <-m.MEvent:
			switch e.Event {
			case memberlist.NodeJoin:
				m.OnAdd(m.toNode(e.Node))
			case memberlist.NodeLeave:
				m.OnRemove(m.toNode(e.Node))
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
			if mr.Token > 0 {
				nodes := m.RingNode(mr.Token, mr.Replicas)
				mr.Async <- nodes
			} else {
				nodes := make([]core.Node, 0)
				for _, n := range m.nodes {
					nodes = append(nodes, n)
				}
				mr.Async <- nodes
			}
		}
	}
}

func (m *MemberListListener) KeyRing(r core.RingRequest) {
	r.Replicas = REPLICA_MAX
	m.MRequest <- r
}

func (m *MemberListListener) HashRing(r core.RingRequest) {
	m.MRequest <- r
	m.UpdateNode(500 * time.Millisecond)
}

func (m *MemberListListener) Get(get core.GetRequest) {
	rq := make(chan []core.Node, 1)
	defer close(rq)
	retry := RetryTrack{Reties: RETRY_MAX}

	for retry.Reties > 0 {
		m.MRequest <- core.RingRequest{Token: m.RingToken(get.Key), Async: rq}
		nodes := <-rq
		ringNode := nodes[0]
		//core.AppLog.Printf("target node %s %s %d\n", ringNode.IP, ringNode.Name, ringNode.RingToken)
		dt, err := m.MemberDataListener.Get(&ringNode, &get)
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
		resp, err := m.MemberDataListener.Set(&ringNode, &set)
		if err != nil {
			retry.Err = err
			retry.Reties--
			continue
		}
		set.Async <- core.Chunk{Remaining: false, Data: []byte(resp.Message)}
		retry.Suc = true
		slaves := nodes[1:]
		for _, slave := range slaves {
			m.MemberDataListener.Set(&slave, &set)
		}
		break
	}
	if retry.Suc {
		return
	}
	core.AppLog.Printf("retry %s, %d", retry.Err.Error(), retry.Reties)
	set.Async <- core.Chunk{Remaining: false, Data: []byte(retry.Err.Error())}
}

func (m *MemberListListener) ShutdownHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	core.AppLog.Println("Signal to exit")
	m.Leave(3 * time.Second)
	time.Sleep(5 * time.Second)
	m.Shutdown()
	signal.Stop(sigs)
	close(sigs)
	os.Exit(0)
}

// delegate
func (m *MemberListListener) NodeMeta(limit int) []byte {
	//core.AppLog.Printf("pull meta data for node update")
	m.ct++
	return fmt.Appendf([]byte{}, "tarantula%d", m.ct)
}

func (m *MemberListListener) NotifyMsg(msg []byte) {
	core.AppLog.Printf("notify msf %s", string(msg))
}

func (m *MemberListListener) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (m *MemberListListener) LocalState(join bool) []byte {
	if join {
		return []byte("mice")
	}
	return []byte("dog")
}
func (m *MemberListListener) MergeRemoteState(buf []byte, join bool) {
	//core.AppLog.Printf("MergeRemoteState %s %v\n", string(buf), join)
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
