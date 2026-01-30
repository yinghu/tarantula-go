package main

import (
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
	MConflict chan core.Node
	MRequest  chan core.RingRequest
	*memberlist.Memberlist
	*MemberHashRing
}

// event dispatch from event delegate
func (m *MemberListListener) Listen() {
	for {
		select {
		case e := <-m.MEvent:
			switch e.Event {
			case memberlist.NodeJoin:
				m.Add(core.Node{Name: e.Node.Name})
			case memberlist.NodeLeave:
				m.Remove(core.Node{Name: e.Node.Name})
			case memberlist.NodeUpdate:
				m.Update(core.Node{Name: e.Node.Name})
			}
			core.AppLog.Printf("Cluster event : %v\n", e)
		case mg := <-m.MMerge:
			core.AppLog.Printf("Merge event %v\n", mg)
		case ma := <-m.MAlive:
			core.AppLog.Printf("Alive event %v\n", ma)
			m.UpdateNode(5 * time.Second)
		case mp := <-m.MPing:
			core.AppLog.Printf("Ping event %v\n", mp)
		case mc := <-m.MConflict:
			core.AppLog.Printf("Conflict event %v\n", mc)
		case mr := <-m.MRequest:
			if mr.Token > 0 {
				node := m.FindNode(mr.Token)
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
	return []byte("tarantula")
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
	m.MPing <- core.Node{Name: other.Name}
}

// merge delegate
func (m *MemberListListener) NotifyMerge(peers []*memberlist.Node) error {
	nodes := make([]core.Node, 0, len(peers))
	for _, n := range peers {
		nodes = append(nodes, core.Node{Name: n.Name})
	}
	m.MMerge <- nodes
	return nil
}

// alive delegate
func (m *MemberListListener) NotifyAlive(peer *memberlist.Node) error {
	m.MAlive <- core.Node{Name: peer.Name}
	return nil
}

// conflict delegate
func (m *MemberListListener) NotifyConflict(existing, other *memberlist.Node) {
	m.MConflict <- core.Node{Name: other.Name}
}
