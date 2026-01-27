package main

import (
	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
)

type MemberlistManager struct {
	Seed []string
	MemberListener
	Ch chan memberlist.NodeEvent
}

func (m *MemberlistManager) Start() error {
	cfg := memberlist.DefaultLANConfig()
	ch := make(chan memberlist.NodeEvent, 10) //HAVE TO BUFFER
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	m.Ch = ch
	cfg.Logger = core.AppLog
	cfg.Events = &cl
	cfg.Delegate = m
	cfg.Ping = m
	list, err := memberlist.Create(cfg)
	if err != nil {
		core.AppLog.Printf("erorr on member create %s\n", err.Error())
		return err
	}
	joined, err := list.Join(m.Seed)
	if err != nil {
		core.AppLog.Printf("erorr on member join %s\n", err.Error())
		return err
	}
	core.AppLog.Printf("joined %d\n", joined)
	return nil
}
