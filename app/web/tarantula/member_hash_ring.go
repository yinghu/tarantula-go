package main

import (
	"fmt"
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
	nodes   []core.Node
	weight  int
	nodeNum int
	WNode   chan<- []core.Node
}

func (m *MemberHashRing) vNode(node core.Node, weight int) core.Node {
	v := core.Node{Name: fmt.Sprintf("%s#%d", node.Name, weight), IP: node.IP, Meta: node.Meta, State: node.State}
	v.RingToken = m.RingToken([]byte(v.Name))
	return v
}

func (m *MemberHashRing) OnAdd(node core.Node) {
	added := make([]core.Node, 0, m.weight)
	for w := range m.weight {
		v := m.vNode(node, w)
		node.RingToken = m.RingToken([]byte(v.Name))
		m.nodes = append(m.nodes, v)
		added = append(added, v)
	}
	slices.SortFunc(m.nodes, cmp)
	m.nodeNum++
	m.WNode <- added
}

func (m *MemberHashRing) OnRemove(node core.Node) {
	removed := make([]core.Node, 0,m.weight)
	m.nodes = slices.DeleteFunc(m.nodes, func(n core.Node) bool {
		if n.IP == node.IP {
			n.State = NODE_STATE_DEAD
			removed = append(removed, n)
			return true
		}
		return false
	})
	slices.SortFunc(m.nodes, cmp)
	m.nodeNum--
	m.WNode <- removed
}

func (m *MemberHashRing) OnUpdate(node core.Node) {

}

func (m *MemberHashRing) OnMerge(nodes []core.Node) {

}

func (m *MemberHashRing) OnLive(node core.Node) {

}

func (m *MemberHashRing) OnPing(node core.Node) {

}

func (m *MemberHashRing) OnConflict(nodes []core.Node) {

}

// hash ring operations
func (m *MemberHashRing) RingToken(key []byte) uint32 {
	return murmur3.Sum32(key)
}

func (m *MemberHashRing) rangeNodeAdded(t uint32) []core.Node {
	target := core.Node{RingToken: t}
	n, exist := slices.BinarySearchFunc(m.nodes, target, func(n1, n2 core.Node) int {
		return cmp(n1, n2)
	})
	if exist {
		if n > 0 {
			return []core.Node{m.nodes[n-1]}
		} else {
			p := len(m.nodes) - 1
			return []core.Node{m.nodes[p]}
		}
	}
	return []core.Node{target}
}
func (m *MemberHashRing) rangeNodeRemoved(t uint32) []core.Node {
	ix := m.indexOf(t)
	return []core.Node{m.nodes[ix]}
}

func (m *MemberHashRing) keyRing(t uint32, relica int) []core.Node {
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

func (m *MemberHashRing) indexOf(t uint32) int {
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
