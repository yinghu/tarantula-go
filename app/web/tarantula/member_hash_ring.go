package main

import "gameclustering.com/internal/core"

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
}

func (m *MemberHashRing) FindNode(t uint32) core.Node {
	n := core.Node{RingToken: t}
	l := 0
	r := len(m.nodes) - 1
	ix := -1
	for l <= r {
		md := l + (r-l)/2
		if n.RingToken < m.nodes[md].RingToken {
			ix = md
			r = md - 1
		} else {
			l = md + 1
		}
	}
	if ix == -1 {
		ix = 0
	}
	n.Name = m.nodes[ix].Name
	return n
}
