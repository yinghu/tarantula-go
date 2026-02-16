package main

import "gameclustering.com/internal/core"

type MemberDataListener struct {
	RNode <-chan core.Node
}

func (m *MemberDataListener) Listen() {
	for n := range m.RNode{
		core.AppLog.Printf("node updated IP : %s NAME : %s RING TOKEN : %d STATE : %d",n.IP,n.Name,n.RingToken,n.State)
	}
}
