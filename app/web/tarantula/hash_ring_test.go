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
		t.Errorf("Total should be rings %d insead of %d\n", total, len(dup))
	}
}

func TestHashRing(t *testing.T) {
	nodes := make([]core.Node, 0, 5)
	for i := range 5 {
		name := fmt.Sprintf("node-%d", i)
		nodes = append(nodes, core.Node{Name: name, RingToken: uint32(i + 100 + i*10)})
	}
	mr := MemberHashRing{nodes: nodes}
	slices.SortFunc(nodes,cmp)
	for _, n := range mr.nodes {
		fmt.Printf("Node ring %d\n", n.RingToken)
	}
	var rt uint32 = 90
	keyNode := core.Node{RingToken: rt}
	fmt.Printf("Key  ring %v\n", keyNode.RingToken)
	pos := slices.IndexFunc(mr.nodes, func(t core.Node) bool {
		if keyNode.RingToken < t.RingToken {
			keyNode.Name = t.Name
			return true
		}
		return false
	})
	fmt.Printf("node pos %d %s\n", pos, keyNode.Name)
	pn := mr.FindNode(rt)
	fmt.Printf("node found %d %s\n", pn.RingToken, pn.Name)
	if keyNode.Name != pn.Name{
		t.Errorf("should be same node %s <> %s",keyNode.Name,pn.Name)
	}
}
