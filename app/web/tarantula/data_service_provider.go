package main

import (
	context "context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/tarantula/protocol"
	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/grpc"
)

const (
	RPC_PORT int = 7001
)

type DataServiceProvider struct {
	protocol.UnimplementedDataServiceServer
	Local  *persistence.BadgerLocal
	RNode  <-chan []core.Node
	RSync  <-chan []byte
	server *grpc.Server
	Mll    MemberListListener

	//write worker chan
	DSet  chan SetData
	DPull chan core.RingSync
	DWait sync.WaitGroup
}

func (c *DataServiceProvider) Get(ctx context.Context, in *protocol.Request) (*protocol.Data, error) {
	data := protocol.Data{}
	getdata := GetData{in}
	k, err := getdata.K()
	if err != nil {
		return &data, err
	}
	ki := getdata.KeyIndex()
	ak, err := ki.K()
	if err != nil {
		return &data, err
	}
	err = c.Local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(ak)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			ki.V(val)
			return nil
		})
		if err != nil {
			return err
		}
		item, err = txn.Get(k)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			data.Value = val
			return nil
		})
	})
	core.AppLog.Debug().Msgf("key index %v", ki)
	return &data, err
}

func (c *DataServiceProvider) Set(ctx context.Context, in *protocol.Data) (*protocol.Response, error) {
	msg := make(chan SetRes, 1)
	defer close(msg)
	setData := SetData{Data: in, Msg: msg}
	//setData.Header.Revision = 109
	//setData.Header.Timestamp = time.Now().UnixMilli()
	//setData.Header.Size = int32(len(setData.Value))
	c.DSet <- setData
	resp := <-msg
	if resp.Err != nil {
		core.AppLog.Debug().Msgf("set error %s", resp.Err.Error())
	}
	return &protocol.Response{Successful: resp.Suc}, resp.Err
}

func (c *DataServiceProvider) Create(ctx context.Context, in *protocol.Data) (*protocol.Response, error) {
	msg := make(chan SetRes, 1)
	defer close(msg)
	setData := SetData{Data: in, Msg: msg}
	c.DSet <- setData
	resp := <-msg
	return &protocol.Response{Successful: resp.Suc}, resp.Err
}

func (c *DataServiceProvider) Update(ctx context.Context, in *protocol.Data) (*protocol.Response, error) {
	msg := make(chan SetRes, 1)
	defer close(msg)
	setData := SetData{Data: in, Msg: msg}
	c.DSet <- setData
	resp := <-msg
	return &protocol.Response{Successful: resp.Suc}, resp.Err
}

func (c *DataServiceProvider) Delete(ctx context.Context, in *protocol.Data) (*protocol.Response, error) {
	msg := make(chan SetRes, 1)
	defer close(msg)
	setData := SetData{Data: in, Msg: msg}
	c.DSet <- setData
	resp := <-msg
	return &protocol.Response{Successful: resp.Suc}, resp.Err
}

func (c *DataServiceProvider) Pull(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.DataBatch]) error {
	//stream.Send()
	return nil
}

func (c *DataServiceProvider) Start() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("%s/%s", homeDir, "tarantula")
	core.AppLog.Printf("creating path %s if not existed", path)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	c.Local = &persistence.BadgerLocal{Path: path, InMemory: false, LogDisabled: false, GcEnabled: true}
	err = c.Local.Open()
	if err != nil {
		panic(err)
	}
	c.DSet = make(chan SetData, NODE_EVENT_BUFFER_SIZE)
	c.DPull = make(chan core.RingSync, NODE_EVENT_BUFFER_SIZE)
	for n := range SET_OPERATOR_NUM {
		go c.runSetData(n)
	}

	tcp, err := net.Listen("tcp", fmt.Sprintf(":%d", RPC_PORT))
	if err != nil {
		panic(err)
	}
	rpc := grpc.NewServer()
	c.server = rpc
	protocol.RegisterDataServiceServer(rpc, c)
	core.AppLog.Printf("local data service provider started on : %s", tcp.Addr().String())
	err = rpc.Serve(tcp)
	if err != nil {
		panic(err)
	}
}

func (m *DataServiceProvider) RingUpdated() {
	running := true
	for running {
		select {
		case nlist := <-m.RNode:
			//update node to pending
			m.Mll.MRequest <- core.RingRequest{Opt: UPDATE_NODE_OPT, Replicas: 0}
			core.AppLog.Debug().Msg("Ring updating")
			ringReqest := core.RingRequest{Address: nlist[0].IP}
			for _, n := range nlist {
				switch n.State {
				case NODE_STATE_SHUTDOWN:
					running = false
				case NODE_STATE_LIVE:
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
			m.Mll.MRequest <- core.RingRequest{Opt: UPDATE_NODE_OPT, Replicas: 1}
			//update node to ready
			core.AppLog.Debug().Msg("Ring updated")
		case sync := <-m.RSync:
			var ds core.RingSync
			err := json.Unmarshal(sync, &ds)
			if err != nil {
				core.AppLog.Warn().Msgf("cannot parse remote data from %s", string(sync))
			} else {
				m.DWait.Add(SET_OPERATOR_NUM)
				for range SET_OPERATOR_NUM {
					m.DSet <- SetData{Opt: SET_OPT_RECOVER}
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
		m.DSet <- SetData{Opt: SET_OPT_CLOSE}
	}
	close(m.DSet)
	close(m.DPull)
	m.server.Stop()
	m.Local.Close()
	core.AppLog.Info().Msg("local data service provider has stopped")
}

// internal operations
func (m *DataServiceProvider) saveKeyIndex(keyIndex *KeyIndex) error {
	key, value, err := keyIndex.Kv()
	if err != nil {
		return err
	}
	if err := m.Local.Db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	}); err != nil {
		return err
	}
	return nil
}

func (m *DataServiceProvider) loadKeyIndex(keyIndex *KeyIndex) error {

	key, err := keyIndex.K()
	if err != nil {
		return err
	}
	return m.Local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return keyIndex.V(val)
		})
	})
}
