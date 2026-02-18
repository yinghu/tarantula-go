package main

import (
	"fmt"
	"slices"
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

func TestHashRing(t *testing.T) {
	nodes := make([]core.Node, 0, 5)
	for i := range 5 {
		name := fmt.Sprintf("node-%d", i)
		nodes = append(nodes, core.Node{Name: name, RingToken: uint32(i + 100 + i*10)})
	}
	mr := MemberHashRing{nodes: nodes}
	slices.SortFunc(nodes, cmp)
	for _, n := range mr.nodes {
		fmt.Printf("Node ring %d", n.RingToken)
	}
	var rt uint32 = 90
	keyNode := core.Node{RingToken: rt}
	fmt.Printf("Key ring %v", keyNode.RingToken)
	pos := slices.IndexFunc(mr.nodes, func(t core.Node) bool {
		if keyNode.RingToken < t.RingToken {
			keyNode.Name = t.Name
			return true
		}
		return false
	})
	fmt.Printf("node pos %d %s", pos, keyNode.Name)
	pn := mr.RingNode(rt, 0)[0]
	fmt.Printf("node found %d %s", pn.RingToken, pn.Name)
	if keyNode.Name != pn.Name {
		t.Errorf("should be same node %s <> %s", keyNode.Name, pn.Name)
	}
	fmt.Println(string(fmt.Appendf([]byte{}, "tarantula%d", 1)))
	n := min(20, 2)
	fmt.Printf("start min number %d", n)
	for n > 0 {
		fmt.Printf("min number %d", n)
		n--
	}
	fmt.Printf("end min number %d", n)
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
	nodes := ring.RingNode(murmur3.Sum32([]byte("adb")), REPLICA_MAX)
	if len(nodes) != 1 {
		t.Errorf("key ring node should 1 %d", len(nodes))
	}
	ring.OnAdd(core.Node{Name: "node-b", IP: "192.168.1.11:6060"})
	nodes = ring.RingNode(murmur3.Sum32([]byte("adb")), REPLICA_MAX)
	if len(nodes) != 2 {
		t.Errorf("key ring node should 2 %d", len(nodes))
	}

	ring.OnAdd(core.Node{Name: "node-c", IP: "192.168.1.12:6060"})
	nodes = ring.RingNode(murmur3.Sum32([]byte("adb")), REPLICA_MAX)
	if len(nodes) != 3 {
		t.Errorf("key ring node should 3 %d", len(nodes))
	}
	ring.OnAdd(core.Node{Name: "node-d", IP: "192.168.1.13:6060"})
	ring.OnAdd(core.Node{Name: "node-e", IP: "192.168.1.14:6060"})

	nodes = ring.RingNode(murmur3.Sum32([]byte("bopaa")), REPLICA_MAX)
	if len(nodes) != 3 {
		t.Errorf("key ring node should 3 %d", len(nodes))
	}
	fmt.Printf("NODES : %v", nodes)

}

func TestHashRingPrefix(t *testing.T) {
	core.CreateTestLog()
	rwNode := make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	ring := MemberHashRing{weight: NODE_WEIGHT, WNode: rwNode}
	ring.OnAdd(core.Node{Name: "node-a", IP: "192.168.1.10:6060"})
	key := []byte("key1")
	hash := ring.RingToken(key)
	nodes := ring.RingNode(hash, REPLICA_MAX)
	core.AppLog.Debug().Msgf("xnode hash %d", nodes[0].RingToken)
	buff := core.NewBuffer(100)
	buff.WriteUInt32(nodes[0].RingToken)
	buff.Write(key)
	buff.Flip()
	data, _ := buff.Read(0)
	resp := core.NewBuffer(100)
	resp.Write(data)
	resp.Flip()
	h, _ := resp.ReadUInt32()
	k, _ := resp.Read(0)
	core.AppLog.Debug().Uint32("h",h).Str("k",string(k)).Send()

}
