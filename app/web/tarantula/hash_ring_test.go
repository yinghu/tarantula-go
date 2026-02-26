package main

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
	"github.com/spaolacci/murmur3"
)

func TestMurmur3(t *testing.T) {
	hash := murmur3.New32()
	dup := make(map[uint32]int)
	i := 0
	total := 100000
	for {
		if i == total {
			break
		}
		fmt.Fprintf(hash, "node%d", i)
		h := hash.Sum32()
		ring := h // % RING_SLOTS >> duplicate ???

		v, e := dup[ring]
		if !e {
			dup[ring] = 1
		} else {
			dup[ring] = v + 1
		}
		i++
	}
	if len(dup) != total {
		t.Errorf("Total should be rings %d insead of %d", total, len(dup))
	}
}

func TestHashRingScale(t *testing.T) {
	core.CreateTestLog()
	rwNode := make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	ring := MemberHashRing{weight: NODE_WEIGHT, WNode: rwNode}

	ring.OnAdd(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	if len(ring.nodes) != 7 {
		t.Errorf("ring node should 7 %d", len(ring.nodes))
	}
	ring.OnRemove(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	if len(ring.nodes) != 0 {
		t.Errorf("ring node should 0 %d", len(ring.nodes))
	}
	//hash := murmur3.New32()
	ring.OnAdd(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	nodes := ring.keyRing(murmur3.Sum32([]byte("adb")), 3)
	if len(nodes) != 1 {
		t.Errorf("key ring node should 1 %d", len(nodes))
	}
	ring.OnAdd(core.Node{Name: "node-b", IP: "192.168.1.11:6060"})
	nodes = ring.keyRing(murmur3.Sum32([]byte("adb")), 3)
	if len(nodes) != 2 {
		t.Errorf("key ring node should 2 %d", len(nodes))
	}

	ring.OnAdd(core.Node{Name: "node-c", IP: "192.168.1.12:6060"})
	nodes = ring.keyRing(murmur3.Sum32([]byte("adb")), 3)
	if len(nodes) != 3 {
		t.Errorf("key ring node should 3 %d", len(nodes))
	}
	ring.OnAdd(core.Node{Name: "node-d", IP: "192.168.1.13:6060"})
	ring.OnAdd(core.Node{Name: "node-e", IP: "192.168.1.14:6060"})

	nodes = ring.keyRing(murmur3.Sum32([]byte("bopaa")), 3)
	if len(nodes) != 3 {
		t.Errorf("key ring node should 3 %d", len(nodes))
	}
	//fmt.Printf("NODES : %v", nodes)

}

func TestHashRingPrefix(t *testing.T) {
	core.CreateTestLog()
	rwNode := make(chan []core.Node, NODE_EVENT_BUFFER_SIZE*100)
	ring := MemberHashRing{weight: NODE_WEIGHT, WNode: rwNode}
	ring.OnAdd(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	ring.OnAdd(core.Node{Name: "node-b", IP: "192.168.1.11:6060"})
	ring.OnAdd(core.Node{Name: "node-c", IP: "192.168.1.12:6060"})
	ring.OnAdd(core.Node{Name: "node-d", IP: "192.168.1.13:6060"})
	ring.OnAdd(core.Node{Name: "node-e", IP: "192.168.1.14:6060"})
	ring.OnAdd(core.Node{Name: "node-f", IP: "192.168.1.15:6060"})
	mindex := make(map[string]uint32)
	key := []byte("key1")
	hash := ring.RingToken(key)
	nodes := ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		mindex[n.Name] = n.RingToken
	}
	ring.OnRemove(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	shared := 0
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0
	ring.OnAdd(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0
	ring.OnAdd(core.Node{Name: "node-g", IP: "192.168.1.7:6060"})
	ring.OnAdd(core.Node{Name: "node-h", IP: "192.168.1.8:6060"})
	ring.OnAdd(core.Node{Name: "node-i", IP: "192.168.1.9:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0
	ring.OnAdd(core.Node{Name: "node-j", IP: "192.168.1.17:6060"})
	ring.OnAdd(core.Node{Name: "node-k", IP: "192.168.1.18:6060"})
	ring.OnAdd(core.Node{Name: "node-l", IP: "192.168.1.19:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0
	ring.OnAdd(core.Node{Name: "node-m", IP: "192.168.1.27:6060"})
	ring.OnAdd(core.Node{Name: "node-n", IP: "192.168.1.28:6060"})
	ring.OnAdd(core.Node{Name: "node-o", IP: "192.168.1.29:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0
	ring.OnAdd(core.Node{Name: "node-q", IP: "192.168.1.37:6060"})
	ring.OnAdd(core.Node{Name: "node-r", IP: "192.168.1.38:6060"})
	ring.OnAdd(core.Node{Name: "node-s", IP: "192.168.1.39:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0
	ring.OnAdd(core.Node{Name: "node-x", IP: "192.168.1.47:6060"})
	ring.OnAdd(core.Node{Name: "node-y", IP: "192.168.1.48:6060"})
	ring.OnAdd(core.Node{Name: "node-z", IP: "192.168.1.49:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0
	ring.OnAdd(core.Node{Name: "node-v", IP: "192.168.1.147:6060"})
	ring.OnAdd(core.Node{Name: "node-t", IP: "192.168.1.148:6060"})
	ring.OnAdd(core.Node{Name: "node-u", IP: "192.168.1.149:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0

	ring.OnAdd(core.Node{Name: "node-aa", IP: "192.168.1.247:6060"})
	ring.OnAdd(core.Node{Name: "node-bb", IP: "192.168.1.248:6060"})
	ring.OnAdd(core.Node{Name: "node-cc", IP: "192.168.1.249:6060"})
	nodes = ring.keyRing(hash, REPLICA_MAX)
	for _, n := range nodes {
		//core.AppLog.Debug().Msgf("hash %d ip %s name %s", n.RingToken, n.IP, n.Name)
		_, ex := mindex[n.Name]
		if ex {
			shared++
		}
	}
	//core.AppLog.Debug().Msgf("shared node number %d", shared)
	shared = 0

	//core.AppLog.Debug().Msgf("total nodes : %d %d", len(ring.nodes)/NODE_WEIGHT, len(ring.nodes))
	//core.AppLog.Debug().Msgf("xnode hash %d", nodes[0].RingToken)
	buff := core.NewBuffer(100)
	buff.WriteUInt32(nodes[0].RingToken)
	buff.WriteUInt32(nodes[0].RingToken)
	buff.WriteUInt32(nodes[0].RingToken)
	buff.Write(key)
	buff.Flip()
	data, _ := buff.Read(0)
	resp := core.NewBuffer(100)
	resp.Write(data)
	resp.Flip()
	//h, _ := resp.ReadUInt32()
	resp.ReadUInt32()
	resp.ReadUInt32()
	//k, _ := resp.Read(0)
	//core.AppLog.Debug().Uint32("h", h).Str("k", string(k)).Send()

}

func h32(n string) uint32 {
	return murmur3.Sum32([]byte(n))
}

func check(t core.Node,rangeRing []core.Node) bool{
	a0 := rangeRing[0]
	a1 := rangeRing[1]
	core.AppLog.Debug().Msgf("pevious : %d <<%d>> next : %d", a0.RingToken, t.RingToken, a1.RingToken)
	order := a0.RingToken > t.RingToken && t.RingToken > a0.RingToken
	return order
}
func TestHashRingBalance(t *testing.T) {
	core.CreateTestLog()
	rwNode := make(chan []core.Node, NODE_EVENT_BUFFER_SIZE*100)
	defer close(rwNode)
	ring := MemberHashRing{weight: NODE_WEIGHT, WNode: rwNode}
	ring.OnAdd(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	for _, n := range ring.nodes {
		core.AppLog.Debug().Msgf("one node ring %v", n)
	}
	for i := range 7 {
		n0 := core.Node{RingToken: h32(fmt.Sprintf("node-a#%d", i))}
		check(n0,ring.rangeNodeAdded(n0.RingToken))
	}

	ring.OnAdd(core.Node{Name: "node-b", IP: "192.168.1.11:6060"})

	for _, n := range ring.nodes {
		core.AppLog.Debug().Msgf("two node ring %v", n)
	}

	for i := range 7 {
		n0 := core.Node{RingToken: h32(fmt.Sprintf("node-a#%d", i))}
		check(n0,ring.rangeNodeAdded(n0.RingToken))	
		n1 := core.Node{RingToken: h32(fmt.Sprintf("node-b#%d", i))}
		check(n1,ring.rangeNodeAdded(n1.RingToken))
	}

	ring.OnAdd(core.Node{Name: "node-c", IP: "192.168.1.12:6060"})
	for _, n := range ring.nodes {
		core.AppLog.Debug().Msgf("three node ring %v", n)
	}
	for i := range 7 {
		n0 := core.Node{RingToken: h32(fmt.Sprintf("node-a#%d", i))}
		check(n0,ring.rangeNodeAdded(n0.RingToken))	
		n1 := core.Node{RingToken: h32(fmt.Sprintf("node-b#%d", i))}
		check(n1,ring.rangeNodeAdded(n1.RingToken))
		n2 := core.Node{RingToken: h32(fmt.Sprintf("node-c#%d", i))}
		check(n2,ring.rangeNodeAdded(n2.RingToken))
	}
	

}
