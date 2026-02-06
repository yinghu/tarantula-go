package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
)

type MemberListListener struct {
	MEvent    chan memberlist.NodeEvent
	MMerge    chan []core.Node
	MAlive    chan core.Node
	MPing     chan core.Node
	MConflict chan []core.Node
	MRequest  chan core.RingRequest
	*memberlist.Memberlist
	*MemberHashRing
	ct int
}

func (m *MemberListListener) toNode(e *memberlist.Node) core.Node {
	return core.Node{Name: e.Name, Meta: string(e.Meta), IP: e.Address(), State: int(e.State)}
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
				node := m.RingNode(mr.Token)
				mr.Async <- []core.Node{node}
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
	m.MRequest <- r
}

func (m *MemberListListener) HashRing(r core.RingRequest) {
	m.MRequest <- r
	m.UpdateNode(500 * time.Millisecond)
}
func (m *MemberListListener) FindValue(r core.ValueRequest) {
	r.Async <- core.Chunk{Remaining: true, Data: []byte("chunk1")}
	r.Async <- core.Chunk{Remaining: false, Data: []byte("chunk2")}
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
	core.AppLog.Printf("pull meta data for node update")
	m.ct++
	return fmt.Appendf([]byte{}, "tarantula%d", m.ct)
}

func (m *MemberListListener) NotifyMsg(msg []byte) {
	core.AppLog.Printf("notify msf %s\n", string(msg))
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
	m.MConflict <- []core.Node{m.toNode(existing),m.toNode(other)}
}
