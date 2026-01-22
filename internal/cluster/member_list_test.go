package cluster

import (
	"fmt"
	"sync"
	"testing"

	"github.com/hashicorp/memberlist"
)

func TestMemberList(t *testing.T) {
	cfg := memberlist.DefaultLANConfig()
	ch := make(chan memberlist.NodeEvent)
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	cfg.Name = "a01"
	cfg.Events = &cl

	fmt.Printf("config %s %d %v\n", cfg.Name, cfg.BindPort, cfg.SecretKey)
	list, err := memberlist.Create(cfg)
	if err != nil {
		fmt.Printf("Erorr %s\n", err.Error())
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(m *memberlist.Memberlist) {
		for e := range ch {
			fmt.Printf("Node event %v\n", e)
		}
	}(list)
	//n, err := list.Join([]string{"192.168.1.11:7946"})
	//if err != nil {
	//panic("Failed to join cluster: " + err.Error())
	//}
	//fmt.Printf("joined : %d\n", n)
	// Ask for members of the cluster
	wg.Wait()
}
