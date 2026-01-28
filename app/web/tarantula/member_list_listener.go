package main

import (
	"time"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
)

type MemberListListener struct {
	Ch chan memberlist.NodeEvent
	*memberlist.Memberlist
}

// event dispatch from event delegate
func (m *MemberListListener) Listen() {
	for e := range m.Ch {
		core.AppLog.Printf("Cluster event : %v %d\n", e, m.NumMembers())
	}
}

func (m *MemberListListener) List() []core.Node{
	nodes := make([]core.Node,0)
	for _, n := range m.Members() {
		node := core.Node{Name: n.Name}
		nodes = append(nodes,node)
	}
	return nodes
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
