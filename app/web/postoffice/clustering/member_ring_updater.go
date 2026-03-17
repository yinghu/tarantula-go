package clustering

import (
	"encoding/json"

	"gameclustering.com/internal/core"
)

type RingUpdate struct {
	State int
	Nodes []core.Node
}

func (m *DataServiceProvider) RingUpdated() {
	running := true
	for running {
		select {
		case ringUpdate := <-m.RNode:
			if ringUpdate.State == NODE_STATE_SHUTDOWN {
				running = false
			} else {
				ringReqest := core.RingRequest{Address: ringUpdate.Nodes[0].IP}
				for _, n := range ringUpdate.Nodes {
					switch n.State {
					case NODE_STATE_LIVE:
						core.AppLog.Debug().Msgf("node added %v", n)

						if !m.Mll.localNode(n) { //skip node initial add call
							rq := make(chan []core.Node, 1)
							m.Mll.rangeRing(core.RingRequest{Token: n.RingToken, Opt: ADD_NODE_OPT, Async: rq})
							ringRange := <-rq
							close(rq)
							if m.Mll.localNode(ringRange[1]) {
								ringReqest.Opt = SYNC_NODE_OPT
								ringReqest.Source.Remote = ringRange[1].RpcEndpoint
								ringReqest.Source.Hashs = append(ringReqest.Source.Hashs, ringRange[0].RingToken)
								core.AppLog.Debug().Msgf("push data key hash >= %d and < %d to remote node %s", ringRange[0].RingToken, n.RingToken, n.IP)
							}
						}

					case NODE_STATE_DEAD:
						if !m.Mll.localNode(n) { //skip node initial add call
							rq := make(chan []core.Node, 1)
							m.Mll.rangeRing(core.RingRequest{Token: n.RingToken, Opt: REMOVE_NODE_OPT, Async: rq})
							ringRange := <-rq
							close(rq)
							if !m.Mll.localNode(ringRange[0]) {
								//pull remote data from >= pre.hash to < added.hash to remote added node
								core.AppLog.Debug().Msgf("take over data key hash >= %d and < %d to remote node %s", ringRange[0].RingToken, n.RingToken, n.IP)
							}
						}
					}
				}
				if ringReqest.Opt == SYNC_NODE_OPT {
					m.Mll.MRequest <- ringReqest
				}
			}
		case sync := <-m.RSync:
			var ds core.RingSync
			err := json.Unmarshal(sync, &ds)
			if err != nil {
				core.AppLog.Warn().Msgf("cannot parse remote data from %s", string(sync))
			} else {
				m.DWait.Add(SET_OPERATOR_NUM)
				for range SET_OPERATOR_NUM {
					m.DSet <- SetData{Opt: core.SET_OPT_RECOVER}
				}
				pz := SET_OPERATOR_NUM
				for _, p := range ds.Hashs {
					ps := core.RingSync{Remote: ds.Remote, Hashs: []uint32{p}}
					m.DPull <- ps
					pz--
				}
				for range pz {
					m.DPull <- core.RingSync{Remote: ds.Remote}
				}
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
