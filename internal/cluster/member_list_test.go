package cluster

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/memberlist"
)

func TestMemberList(t *testing.T) {
	cfg := memberlist.DefaultLANConfig()
	cfg.Name = "a02"
	fmt.Printf("config %s %d %v\n", cfg.Name, cfg.BindPort, cfg.SecretKey)
	list, err := memberlist.Create(cfg)
	if err != nil {
		fmt.Printf("Erorr %s\n", err.Error())
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(m *memberlist.Memberlist) {
		tk := time.NewTicker(5 * time.Second)
		for {
			t := <-tk.C
			for _, member := range m.Members() {

				fmt.Printf("Member: %s %s %s %v\n", member.Name, member.Addr, m.LocalNode().Name, t)
			}
		}
		//wg.Wait()
	}(list)
	//n, err := list.Join([]string{"192.168.1.11:7946"})
	//if err != nil {
		//panic("Failed to join cluster: " + err.Error())
	//}
	//fmt.Printf("joined : %d\n", n)
	// Ask for members of the cluster
	wg.Wait()
}
