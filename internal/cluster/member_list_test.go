package cluster

import (
	"fmt"
	"sync"
	"testing"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
)

func TestMemberList(t *testing.T) {
	core.CreateTestLog()
	cfg := memberlist.DefaultLANConfig()
	ch := make(chan memberlist.NodeEvent, 10) //HAVE TO BUFFER
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	cfg.Logger = core.AppLog
	cfg.Events = &cl
	fmt.Printf("config %s %d %v %s\n", cfg.Name, cfg.BindPort, cfg.SecretKey, cfg.Label)
	list, err := memberlist.Create(cfg)
	if err != nil {
		fmt.Printf("Erorr %s\n", err.Error())
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(m *memberlist.Memberlist) {
		for {
			e := <-ch
			fmt.Printf("Node event %v %d\n", e, m.NumMembers())
		}
	}(list)
	//fmt.Printf("joining to ")
	n, err := list.Join([]string{"192.168.1.11:7946", "192.168.1.6:7946"})
	if err != nil {
		panic("Failed to join cluster: " + err.Error())
	}
	fmt.Printf("joined : %d\n", n)
	// Ask for members of the cluster
	wg.Wait()
}
