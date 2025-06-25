package cluster

import (
	"cmp"
	"fmt"
	"slices"
	"testing"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

func TestPartition(t *testing.T) {
	nodes := []string{"a11", "a02", "a03", "a04", "a05", "a06", "a07", "a08", "a09", "a10", "a01", "a12"}
	slices.Sort(nodes)
	mp := make(map[int]string)
	sz := len(nodes)
	for p := range 31 {
		i := p % sz
		//fmt.Printf("Partition %d %s %d\n", i, nodes[i], p)
		mp[p] = nodes[i]
	}
}

func TestNode(t *testing.T) {
	nodes := []core.Node{{Name: "a05"}, {Name: "a04"}, {Name: "a02"}, {Name: "a01"}}
	slices.SortFunc(nodes, func(a, b core.Node) int {
		return cmp.Compare(a.Name, b.Name)
	})
	for n := range nodes {
		fmt.Printf("Node : %s %d\n", nodes[n].Name, n)
	}
}

func TestCluster(t *testing.T) {
	lc := core.Node{Name: "a01", HttpEndpoint: "http://192.168.1.11:8080", TcpEndpoint: "tcp://192.168.1.11:5000"}
	c := NewEtc("tarantula", []string{"192.168.1.7:2379"}, LocalNode{Node: lc})
	tk := time.NewTimer(20 * time.Second)
	go func() {
		<-tk.C
		c.Quit <- true
	}()
	err := c.Join()
	if err != nil {
		t.Errorf("Service error %s", err.Error())
	}
}

func TestTransaction(t *testing.T) {
	lc := core.Node{Name: "a01", HttpEndpoint: "http://192.168.1.11:8080", TcpEndpoint: "tcp://192.168.1.11:5000"}
	c := NewEtc("tarantula", []string{"192.168.1.7:2379"}, LocalNode{Node: lc})
	c.Atomic("", func(c core.Ctx) error {
		k := "jwkkey"
		v := util.KeyToBase64(util.Key(32))
		err := c.Put(k, v)
		if err != nil {
			t.Errorf("Put error %s", err.Error())
			return err
		}
		gv, err := c.Get(k)
		if err != nil {
			t.Errorf("Get error %s", err.Error())
			return err
		}
		if gv != v {
			t.Errorf("value not same error %s", gv)
		}
		err = c.Del(k)
		if err != nil {
			t.Errorf("Del error %s", err.Error())
		}
		dv, err := c.Get(k)
		if err == nil {
			t.Errorf("Get error %s", dv)
		}
		return nil
	})
}
