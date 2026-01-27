package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/memberlist"
)

type MemberListener struct {
		
}

func (m *MemberListener) NodeMeta(limit int) []byte {
	fmt.Printf("node meta %d\n", limit)
	return []byte("tarantula")
}

func (m *MemberListener) NotifyMsg(msg []byte) {
	fmt.Printf("notify msf %s\n", string(msg))
}

func (m *MemberListener) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (m *MemberListener) LocalState(join bool) []byte {
	if join {
		return []byte("mice")
	}
	return []byte("dog")
}
func (m *MemberListener) MergeRemoteState(buf []byte, join bool) {
	fmt.Printf("MergeRemoteState %s %v\n", string(buf), join)
}

func (m *MemberListener) AckPayload() []byte {
	return []byte("cat")
}

func (m *MemberListener) NotifyPingComplete(other *memberlist.Node, rtt time.Duration, payload []byte) {
	fmt.Printf("ping :%v %s\n", other, string(payload))
}
