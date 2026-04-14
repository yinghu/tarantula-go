package clustering

import (
	"fmt"
	"time"

	"gameclustering.com/internal/core"
	"github.com/hashicorp/memberlist"
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

	StoreDir string
	Binding  string
	//Seq      core.Sequence
}

func (m *MemberlistManager) Start(meta []byte, seq core.Sequence) error {
	m.MemberHashRing = &MemberHashRing{weight: NODE_WEIGHT}
	m.nodes = make([]core.Node, 0)
	cfg := memberlist.DefaultLANConfig()
	cfg.Name = m.Binding
	ch := make(chan memberlist.NodeEvent, NODE_EVENT_BUFFER_SIZE) //HAVE TO BUFFER
	cl := memberlist.ChannelEventDelegate{Ch: ch}
	m.MEvent = ch
	m.MMerge = make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MAlive = make(chan core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MPing = make(chan core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MConflict = make(chan []core.Node, NODE_EVENT_BUFFER_SIZE)
	m.MRequest = make(chan core.RingRequest, NODE_EVENT_BUFFER_SIZE)
	rwNode := make(chan RingUpdate, NODE_EVENT_BUFFER_SIZE)
	rwSync := make(chan []byte, NODE_EVENT_BUFFER_SIZE)
	m.WNode = rwNode
	m.MSync = rwSync
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
	m.meta = meta
	m.Memberlist = list
	go m.Listen()
	m.DataServiceProvider = &DataServiceProvider{RNode: rwNode, RSync: rwSync, seq: seq}
	m.DataServiceProvider.rpcEndpoint = fmt.Sprintf("%s:%d", string(m.LocalNode().Addr.String()), core.RPC_PORT)
	m.Mll = &m.MemberListListener
	m.Mll.DWait.Add(1)
	go m.DataServiceProvider.Start(m.StoreDir)
	m.Mll.DWait.Wait()
	list.UpdateNode(time.Second * 5)
	go m.RingUpdated()
	joined, err := list.Join(m.Seed)
	if err != nil {
		core.AppLog.Printf("erorr on member join %s", err.Error())
		return err
	}
	core.AppLog.Printf("total nodes have joined %d on local node  %s", joined, m.DataServiceProvider.rpcEndpoint)
	return nil
}

func (m *MemberlistManager) ShutdownHook() {
	core.AppLog.Info().Msg("running shut down hook ...")
	m.running = false
	m.Leave(3 * time.Second)
	time.Sleep(3 * time.Second)
	m.Shutdown()
	core.AppLog.Info().Msg("closing resouces ...")
	m.MRequest <- core.RingRequest{Opt: CLOSE_RING_OPT}
	m.WNode <- RingUpdate{State: NODE_STATE_SHUTDOWN}
	time.Sleep(3 * time.Second)
	close(m.MEvent)
	close(m.MAlive)
	close(m.MPing)
	close(m.MMerge)
	close(m.MConflict)
	close(m.MRequest)
	close(m.WNode)
	close(m.MSync)
	core.AppLog.Info().Msg("shut down has done successfully.")
}
