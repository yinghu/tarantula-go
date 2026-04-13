package clustering

import (
	"encoding/json"
	"slices"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

type RingUpdate struct {
	State int
	Nodes []core.Node
}

const (
	RECEIVER_START  uint32 = 1
	TOPIC_REGISTER  uint32 = 2
	RECEIVER_REMOVE uint32 = 3
	RECEIVER_END    uint32 = 4
)

type ReceiverAsync struct {
	Rev chan *protocol.Mail
	Q   chan string
}

type TopicRequest struct {
	Opt    uint32
	NodeId string
	Tag    string
	Name   string

	Async chan ReceiverAsync
	Subs  chan []core.Subscription
}

func (m *DataServiceProvider) balanceOnNodeAdded(added RingUpdate) {

	if m.backRing.nodeNum == 0 {
		m.backRing.nodes = append(m.backRing.nodes, added.Nodes...)
		slices.SortFunc(m.backRing.nodes, cmp)
		m.backRing.nodeNum++
		return
	}
	slices.SortFunc(added.Nodes, cmp)
	ringSync := core.RingSync{Ranges: make([]core.RingRange, 0)}
	for _, n := range added.Nodes {
		if !m.Mll.localNode(n) { //skip node initial add call
			ringRange := m.backRing.rangeOfRing(n.RingToken)
			if m.Mll.localNode(ringRange[1]) {
				ringSync.Remote = ringRange[1].RpcEndpoint
				ringSync.Ranges = append(ringSync.Ranges, core.RingRange{From: ringRange[0].RingToken, To: n.RingToken})
				core.AppLog.Debug().Msgf("push data key hash >= %d and < %d to remote node %s", ringRange[0].RingToken, n.RingToken, n.IP)
			}
		}
		m.backRing.nodes = append(m.backRing.nodes, n)
		slices.SortFunc(m.backRing.nodes, cmp)
	}
	m.backRing.nodeNum++
	if len(ringSync.Ranges) == 0 {
		return
	}
	m.Mll.MRequest <- core.RingRequest{Source: ringSync, Opt: SYNC_NODE_OPT, Address: added.Nodes[0].IP}
	m.subscriptions.lookup(func(sub core.Subscription) {
		core.AppLog.Debug().Msgf("sub %v > %s", sub, m.rpcEndpoint)
		if sub.Endpoint == m.rpcEndpoint {
			m.Mll.MRequest <- core.RingRequest{Opt: SYNC_SUB_OPT, Address: added.Nodes[0].IP, Source: core.RingSync{Sub: sub}}
		}
	})
}

func (m *DataServiceProvider) balanceOnNodeRemoved(removed RingUpdate) {

	for _, n := range removed.Nodes {
		m.backRing.nodes = slices.DeleteFunc(m.backRing.nodes, func(d core.Node) bool {
			return d.IP == n.IP
		})
	}
	slices.SortFunc(m.backRing.nodes, cmp)
	m.backRing.nodeNum--
}

func (m *DataServiceProvider) registerSubscription(sub core.Subscription) {
	core.AppLog.Debug().Msgf("subscription %v", sub)
	listener := m.listeners[sub.NodeId]
	core.AppLog.Debug().Msgf("lis %v",listener)
	if sub.Deleting {
		m.subscriptions.del(sub)
	} else {
		m.subscriptions.add(sub)
	}
}

func (m *DataServiceProvider) RingUpdated() {
	running := true
	for running {
		select {
		case ringUpdate := <-m.RNode:
			switch ringUpdate.State {
			case NODE_STATE_SHUTDOWN:
				running = false
			case NODE_STATE_LIVE:
				m.balanceOnNodeAdded(ringUpdate)
			case NODE_STATE_DEAD:
				m.balanceOnNodeRemoved(ringUpdate)
			}

		case sync := <-m.RSync:
			var ds core.RingSync
			err := json.Unmarshal(sync, &ds)
			if err != nil {
				core.AppLog.Warn().Msgf("cannot parse remote data from %s", string(sync))
			} else {
				if len(ds.Ranges) > 0 {
					m.DWait.Add(SET_OPERATOR_NUM)
					for range SET_OPERATOR_NUM {
						m.DSet <- SetData{Opt: core.SET_OPT_RECOVER}
					}
					pz := SET_OPERATOR_NUM
					for _, p := range ds.Ranges {
						ps := core.RingSync{Remote: ds.Remote, Ranges: []core.RingRange{p}}
						m.DPull <- ps
						pz--
					}
					for range pz {
						m.DPull <- core.RingSync{Remote: ds.Remote, Ranges: []core.RingRange{}}
					}
				} else {
					m.registerSubscription(ds.Sub)
				}
			}
		case req := <-m.DRequest:
			switch req.Opt {
			case RECEIVER_START:
				rev, ok := m.listeners[req.Name]
				if !ok {
					rev = ReceiverAsync{Rev: make(chan *protocol.Mail, NODE_EVENT_BUFFER_SIZE), Q: make(chan string, 2)}
					m.listeners[req.Name] = rev
				}
				req.Async <- rev
			case RECEIVER_REMOVE:
				delete(m.listeners, req.Name)
				core.AppLog.Debug().Msgf("listener removed %s", req.Name)
			case TOPIC_REGISTER:
				req.Subs <- m.subscriptions.topic(req)
			case RECEIVER_END:
				rev, ok := m.listeners[req.Name]
				if ok {
					rev.Q <- req.Name
				}
			}
		case msg := <-m.DMessager:

			for _, ch := range m.listeners {
				ch.Rev <- msg
			}

		}

	}
	//shutdown server
	for range SET_OPERATOR_NUM {
		m.DSet <- SetData{Opt: core.SET_OPT_CLOSE}
	}
	close(m.DSet)
	close(m.DPull)
	m.server.Stop()
	m.Local.Close()
	core.AppLog.Info().Msg("local data service provider has stopped")
}
