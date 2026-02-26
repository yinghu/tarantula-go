package main

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/hashicorp/memberlist"
	"github.com/spaolacci/murmur3"
)

const (
	NODE_EVENT_BUFFER_SIZE int = 16
	NODE_WEIGHT            int = 7 //virtual nodes per ip node
	REPLICA_MAX            int = 7
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
	cfg.Events = &cl
	cfg.Delegate = m
	cfg.Ping = m
	cfg.Merge = m
	cfg.Alive = m
	cfg.Conflict = m
	cfg.LogOutput = core.AppLog
	list, err := memberlist.Create(cfg)
	if err != nil {
		core.AppLog.Printf("erorr on member create %s", err.Error())
		return err
	}
	m.Memberlist = list
	m.LocalNode().Meta = []byte("pending")
	nodeId := murmur3.Sum64(m.LocalNode().Addr)
	core.AppLog.Printf("node id %d %d", nodeId, int64(nodeId)%1024)
	m.snowflake = util.NewSnowflake(int64(nodeId)%1024, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	go m.Listen()
	m.DataServiceProvider = &DataServiceProvider{RNode: rwNode}
	m.Mll = m.MemberListListener
	go m.DataServiceProvider.Start()
	go m.RingUpdated()
	joined, err := list.Join(m.Seed)
	if err != nil {
		core.AppLog.Printf("erorr on member join %s", err.Error())
		return err
	}
	core.AppLog.Printf("total nodes have joined %d", joined)
	return nil
}

func (m *MemberlistManager) Id() (int64, error) {
	return m.snowflake.Id()
}
func (m *MemberlistManager) Parse(snowflakeId int64) (int64, int64, int64) {
	return m.snowflake.Parse(snowflakeId)
}

func (m *MemberlistManager) ShutdownHook() {
	core.AppLog.Info().Msg("running shut down hook ...")
	m.Leave(3 * time.Second)
	time.Sleep(3 * time.Second)
	m.Shutdown()
	core.AppLog.Info().Msg("closing resouces ...")
	m.MRequest <- core.RingRequest{Opt: CLOSE_RING_OPT}
	stopNode := []core.Node{{State: NODE_STATE_SHUTDOWN}}
	m.WNode <- stopNode
	time.Sleep(3 * time.Second)
	close(m.MEvent)
	close(m.MAlive)
	close(m.MPing)
	close(m.MMerge)
	close(m.MConflict)
	close(m.MRequest)
	close(m.WNode)
	core.AppLog.Info().Msg("shut down has done successfully.")
}
