package clustering

import (
	"slices"

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
	ix := m.indexOf(t)
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

func (m *NodeRing) rangeNodeAdded(t uint32) []core.Node {
	target := core.Node{RingToken: t}
	n, exist := slices.BinarySearchFunc(m.nodes, target, func(n1, n2 core.Node) int {
		return cmp(n1, n2)
	})
	if exist {
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
	core.AppLog.Warn().Msgf("should be never happening %d %v", n, exist)
	//should not be happening
	return []core.Node{}
}

func (m *NodeRing) rangeNodeRemoved(t uint32) []core.Node {
	n := m.indexOf(t)
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

func (m *NodeRing) indexOf(t uint32) int {
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
