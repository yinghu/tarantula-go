package main

import (
	"slices"

	"gameclustering.com/internal/core"
	"github.com/spaolacci/murmur3"
)

func cmp(n1, n2 core.Node) int {
	if n1.RingToken > n2.RingToken {
		return 1
	}
	if n1.RingToken < n2.RingToken {
		return -1
	}
	return 0
}

type MemberHashRing struct {
	nodes []core.Node
}

func (m *MemberHashRing) Add(node core.Node) {
	core.AppLog.Printf("ADD NODE %v\n", node)
	sz := len(m.nodes)
	node.RingToken = m.RingToken([]byte(node.Name))
	if sz == 0 {
		m.nodes = append(m.nodes, node)
		return
	}
	pending := make([]core.Node, 0, sz+1)
	copy(pending, m.nodes)
	pending = append(pending, node)
	m.nodes = pending
	slices.SortFunc(m.nodes, cmp)
}

func (m *MemberHashRing) Remove(node core.Node) {
	core.AppLog.Printf("REMOVE NODE %v\n", node)
}

func (m *MemberHashRing) Update(node core.Node) {
	core.AppLog.Printf("UPDATE NODE %v\n", node)
}
func (m *MemberHashRing) RingToken(key []byte) uint32 {
	return murmur3.Sum32(key)
}
func (m *MemberHashRing) FindNode(t uint32) core.Node {
	l := 0
	r := len(m.nodes) - 1
	ix := -1
	for l <= r {
		md := l + (r-l)/2
		if t < m.nodes[md].RingToken {
			ix = md
			r = md - 1
		} else {
			l = md + 1
		}
	}
	if ix == -1 {
		ix = 0
	}
	return m.nodes[ix]
}
