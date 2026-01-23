package cluster

import (
	"fmt"
	"sync"
	"testing"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
)

type MockDelegate struct {
}

func (M *MockDelegate) NodeMeta(limit int) []byte {
	fmt.Printf("node meta %d\n", limit)
	return []byte("tarantula")
}

func (M *MockDelegate) NotifyMsg(msg []byte) {
	fmt.Printf("notify msf %s\n", string(msg))
}

func (M *MockDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	//fmt.Printf("broadcast %d %d\n", overhead, limit)
	return nil
}

func (M *MockDelegate) LocalState(join bool) []byte {
	if join{
		fmt.Printf("LocalState %v\n", join)
		return []byte("mice")
	}
	return nil

}
func (M *MockDelegate) MergeRemoteState(buf []byte, join bool) {
	fmt.Printf("MergeRemoteState %s %v\n", string(buf), join)
}
func TestMemberList(t *testing.T) {
	core.CreateTestLog()
	cfg := memberlist.DefaultLANConfig()
	ch := make(chan memberlist.NodeEvent, 10) //HAVE TO BUFFER
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	cfg.Logger = core.AppLog
	cfg.Events = &cl
	cfg.Delegate = &MockDelegate{}

	fmt.Printf("config %s %d %v %s\n", cfg.Name, cfg.BindPort, cfg.SecretKey, cfg.Label)
	list, err := memberlist.Create(cfg)
	if err != nil {
		fmt.Printf("Erorr %s\n", err.Error())
		return
	}
	//list.LocalNode().Meta = []byte("tarantula")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(m *memberlist.Memberlist) {
		for {
			e := <-ch
			fmt.Printf("Node event %v %d %s\n", e, m.NumMembers(), string(m.LocalNode().Meta))
		}
	}(list)
	n, err := list.Join([]string{"192.168.1.11:7946", "192.168.1.6:7946"})
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}
	fmt.Printf("joined : %d\n", n)
	list.SendReliable(list.LocalNode(), []byte("hello"))
	list.SendBestEffort(list.LocalNode(), []byte("udp"))
	//list.Leave(5)
	// Ask for members of the cluster
	wg.Wait()
}
