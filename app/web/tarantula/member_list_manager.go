package main

import (
	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
)

const (
	NODE_EVENT_BUFFER_SIZE int = 16
)

type MemberlistManager struct {
	Seed []string
	MemberListListener
}

func (m *MemberlistManager) Start() error {
	m.MemberHashRing = &MemberHashRing{nodes: make([]core.Node, 0)}
	cfg := memberlist.DefaultLANConfig()
	ch := make(chan memberlist.NodeEvent, NODE_EVENT_BUFFER_SIZE) //HAVE TO BUFFER
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	m.MEvent = ch
	m.MMerge = make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MAlive = make(chan core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MPing = make(chan core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MConflict = make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MRequest = make(chan core.RingRequest, NODE_EVENT_BUFFER_SIZE)
	cfg.Logger = core.AppLog
	cfg.Events = &cl
	cfg.Delegate = m
	cfg.Ping = m
	cfg.Merge = m
	cfg.Alive = m
	cfg.Conflict = m
	list, err := memberlist.Create(cfg)
	if err != nil {
		core.AppLog.Printf("erorr on member create %s\n", err.Error())
		return err
	}
	m.Memberlist = list
	go m.Listen()
	joined, err := list.Join(m.Seed)
	if err != nil {
		core.AppLog.Printf("erorr on member join %s\n", err.Error())
		return err
	}
	core.AppLog.Printf("joined %d\n", joined)
	return nil
}
