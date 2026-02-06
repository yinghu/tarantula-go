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

func (m *MemberHashRing) OnAdd(node core.Node) {
	core.AppLog.Printf("ADD NODE %v\n", node)
	node.RingToken = m.RingToken([]byte(node.Name))
	m.nodes = append(m.nodes, node)
	slices.SortFunc(m.nodes, cmp)
}

func (m *MemberHashRing) OnRemove(node core.Node) {
	core.AppLog.Printf("REMOVE NODE %v\n", node)
	m.nodes = slices.DeleteFunc(m.nodes, func(n core.Node) bool {
		return n.Name == node.Name
	})
	slices.SortFunc(m.nodes, cmp)
}

func (m *MemberHashRing) OnUpdate(node core.Node) {
	core.AppLog.Printf("UPDATE NODE %v\n", node)
}

func (m *MemberHashRing) OnMerge(nodes []core.Node) {
	core.AppLog.Printf("MERGE NODE %v\n", nodes)
}

func (m *MemberHashRing) OnLive(node core.Node) {
	core.AppLog.Printf("LIVE NODE %v\n", node)
}

func (m *MemberHashRing) OnPing(node core.Node) {
	core.AppLog.Printf("PING NODE %v\n", node)
	
}

func (m *MemberHashRing) OnConflict(nodes []core.Node) {
	core.AppLog.Printf("CONFLICT NODE %v\n", nodes)
}


//hash ring operations 
func (m *MemberHashRing) RingToken(key []byte) uint32 {
	return murmur3.Sum32(key)
}

func (m *MemberHashRing) RingNode(t uint32) core.Node {
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
