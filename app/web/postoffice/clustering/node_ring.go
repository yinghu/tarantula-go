package clustering

import (
	"gameclustering.com/internal/core"
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

type NodeRing struct {
	nodes   []core.Node
	nodeNum int
}

func (m *NodeRing) keyRing(t uint32, relica int) []core.Node {
	ix := m.ownerOf(t)
	if relica == 0 || m.nodeNum == 1 {
		return []core.Node{m.nodes[ix]}
	}
	syncNum := min(m.nodeNum, relica) - 1
	syncNodes := make([]core.Node, 0, syncNum)
	sz := len(m.nodes)
	syncNodes = append(syncNodes, m.nodes[ix])
	ix++
	for syncNum > 0 {
		if ix == sz {
			ix = 0
		}
		p := m.nodes[ix]
		dup := false
		for _, nd := range syncNodes {
			if p.IP == nd.IP {
				dup = true
				break
			}
		}
		if !dup {
			syncNum--
			syncNodes = append(syncNodes, p)
		}
		ix++
	}
	return syncNodes
}
func (m *NodeRing) rangeOfRing(t uint32) []core.Node {
	ix := m.ownerOf(t)
	var px int
	if ix == 0 {
		px = len(m.nodes) - 1
	} else {
		px = ix - 1
	}
	return []core.Node{m.nodes[px], m.nodes[ix]}
}

func (m *NodeRing) rangeNodeRemoved(t uint32) []core.Node {
	n := m.ownerOf(t)
	end := len(m.nodes) - 1
	switch n {
	case 0:
		return []core.Node{m.nodes[end], m.nodes[1]}
	case end:
		return []core.Node{m.nodes[0], m.nodes[end-1]}
	default:
		return []core.Node{m.nodes[n-1], m.nodes[n+1]}
	}
}

func (m *NodeRing) ownerOf(t uint32) int {
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
	return ix
}
