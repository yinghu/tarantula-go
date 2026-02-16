package main

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/hashicorp/memberlist"
	"github.com/spaolacci/murmur3"
)

const (
	NODE_EVENT_BUFFER_SIZE int = 16
	NODE_WEIGHT            int = 7 //virtual nodes per ip node
	REPLICA_MAX            int = 3
	RETRY_MAX              int = 3
)

type MemberlistManager struct {
	Seed []string
	MemberListListener
	snowflake util.Snowflake
}

func (m *MemberlistManager) Start() error {
	m.MemberHashRing = &MemberHashRing{nodes: make([]core.Node, 0), weight: NODE_WEIGHT}

	cfg := memberlist.DefaultLANConfig()
	ch := make(chan memberlist.NodeEvent, NODE_EVENT_BUFFER_SIZE) //HAVE TO BUFFER
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	m.MEvent = ch
	m.MMerge = make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MAlive = make(chan core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MPing = make(chan core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MConflict = make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MRequest = make(chan core.RingRequest, NODE_EVENT_BUFFER_SIZE)
	rwNode := make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	m.WNode = rwNode
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
	nodeId := murmur3.Sum64(m.LocalNode().Addr)
	core.AppLog.Printf("node id %d %d\n", nodeId, int64(nodeId)%1024)
	m.snowflake = util.NewSnowflake(int64(nodeId)%1024, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	go m.Listen()
	mdl := MemberDataListener{RNode: rwNode}
	go mdl.Listen()
	joined, err := list.Join(m.Seed)
	if err != nil {
		core.AppLog.Printf("erorr on member join %s\n", err.Error())
		return err
	}
	core.AppLog.Printf("joined %d\n", joined)
	return nil
}

func (m *MemberlistManager) Id() (int64, error) {
	return m.snowflake.Id()
}
func (m *MemberlistManager) Parse(snowflakeId int64) (int64, int64, int64) {
	return m.snowflake.Parse(snowflakeId)
}
