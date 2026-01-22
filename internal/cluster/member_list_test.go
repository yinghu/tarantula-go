package cluster

import (
	"fmt"
	"sync"
	"testing"

	"github.com/hashicorp/memberlist"
)

func TestMemberList(t *testing.T) {
	cfg := memberlist.DefaultLANConfig()
	ch := make(chan memberlist.NodeEvent,10)
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	cfg.Name = "a03"
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
		for {
			e := <-ch
			fmt.Printf("Node event %v\n", e)
		}
	}(list)
	//fmt.Printf("joining to ")
	//n, err := list.Join([]string{"localhost:7946"})
	//if err != nil {
		//panic("Failed to join cluster: " + err.Error())
	//}
	//fmt.Printf("joined : %d\n", n)
	// Ask for members of the cluster
	wg.Wait()
}
