package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
	"github.com/spaolacci/murmur3"
)

const (
	RING_SLOTS int = 271
)

type MemberListListener struct {
	Ch chan memberlist.NodeEvent
	*memberlist.Memberlist
	ringSlots int
	
}

// event dispatch from event delegate
func (m *MemberListListener) Listen() {
	hash := murmur3.New64()
	for e := range m.Ch {
		hash.Write([]byte(e.Node.Name))
		hd := hash.Sum64()
		core.AppLog.Printf("Cluster event : %v %d\n", e, hd%1024)
	}
}

func (m *MemberListListener) List() []core.Node {
	nodes := make([]core.Node, 0)
	for _, n := range m.Members() {
		node := core.Node{Name: n.Name}
		nodes = append(nodes, node)
	}
	return nodes
}

func (m *MemberListListener) ShutdownHook() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	core.AppLog.Println("Signal to exit")
	m.Leave(3 * time.Second)
	m.Shutdown()
	signal.Stop(sigs)
	close(sigs)
	os.Exit(0)
}

// delegate
func (m *MemberListListener) NodeMeta(limit int) []byte {

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
	core.AppLog.Printf("MergeRemoteState %s %v\n", string(buf), join)
}

// ping delegate
func (m *MemberListListener) AckPayload() []byte {
	return []byte("cat")
}

func (m *MemberListListener) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	core.AppLog.Printf("ping :%v %s\n", other, string(payload))
}

// merge delegate
func (m *MemberListListener) NotifyMerge(peers []*memberlist.Node) error {
	core.AppLog.Printf("merge :%v\n", peers)
	return nil
}

// alive delegate
func (m *MemberListListener) NotifyAlive(peer *memberlist.Node) error {
	core.AppLog.Printf("alive :%v\n", peer)
	return nil
}

// conflict delegate
func (m *MemberListListener) NotifyConflict(existing, other *memberlist.Node) {
	core.AppLog.Printf("conflict node :%v %v\n", existing, other)
}
