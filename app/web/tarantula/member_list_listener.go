package main

import (
	"fmt"
	"time"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
)

type MemberListListener struct {
	Ch chan memberlist.NodeEvent
}

func (m *MemberListListener) Listen() {
	for e := range m.Ch {
		core.AppLog.Printf("Cluster event : %v\n", e)
	}
}

func (m *MemberListListener) NodeMeta(limit int) []byte {
	fmt.Printf("node meta %d\n", limit)
	return []byte("tarantula")
}

func (m *MemberListListener) NotifyMsg(msg []byte) {
	fmt.Printf("notify msf %s\n", string(msg))
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
	fmt.Printf("MergeRemoteState %s %v\n", string(buf), join)
}

func (m *MemberListListener) AckPayload() []byte {
	return []byte("cat")
}

func (m *MemberListListener) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	fmt.Printf("ping :%v %s\n", other, string(payload))
}
