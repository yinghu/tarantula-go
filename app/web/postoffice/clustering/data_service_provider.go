package clustering

import (
	context "context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/grpc"
)

type DataServiceProvider struct {
	protocol.UnimplementedDataServiceServer
	protocol.UnimplementedPostofficeServiceServer
	Local  *persistence.BadgerLocal
	RNode  <-chan []core.Node
	RSync  <-chan []byte
	server *grpc.Server
	Mll    *MemberListListener

	//write worker chan
	DSet  chan SetData
	DPull chan core.RingSync
	DWait sync.WaitGroup
}

func (c *DataServiceProvider) Get(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {
	q := event.CreateQuery(request.Query.Id)
	buff := core.NewBuffer(100)
	buff.Write(request.Query.Criteria)
	buff.Flip()
	q.QRead(buff)
	buff.Clear()
	buff.WriteInt32(core.EVENT_FACTORY_ID)
	buff.WriteInt32(event.MESSAGE_CID)
	buff.Flip()
	px, _ := buff.Read(0)
	p := px
	core.AppLog.Debug().Msgf("query : %v", q)
	c.Local.Db.View(func(txn *badger.Txn) error {
		op := badger.IteratorOptions{PrefetchSize: 100, PrefetchValues: false, Reverse: false}
		it := txn.NewIterator(op)
		defer it.Close()

		for it.Seek(p); it.ValidForPrefix(p); it.Next() {
			p = px
			item := it.Item()
			core.AppLog.Debug().Msgf("msg k : %s", string(item.Key()))
			item.Value(func(val []byte) error {
				core.AppLog.Debug().Msgf("msg v : %s", string(val))
				return nil
			})
		}
		return nil
	})

	stream.Send(&protocol.Response{Successful: true, Message: "hello1"})
	stream.Send(&protocol.Response{Successful: true, Message: "hello2"})
	stream.Send(&protocol.Response{Successful: true, Message: "hello3"})
	stream.Send(&protocol.Response{Successful: true, Message: "hello4"})
	stream.Send(&protocol.Response{Successful: true, Message: "hello5"})
	stream.Send(&protocol.Response{Successful: true, Message: "hello6"})
	return nil
}

func (c *DataServiceProvider) Reset(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	msg := make(chan *protocol.Response, 1)
	defer close(msg)
	setData := SetData{Data: in.Data, Resp: msg}
	c.DSet <- setData
	resp := <-msg
	return resp, nil
}

func (c *DataServiceProvider) Create(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	msg := make(chan *protocol.Response, 1)
	defer close(msg)
	setData := SetData{Opt: in.Opt, Data: in.Data, Resp: msg}
	core.AppLog.Debug().Msgf("creating data %d %d", setData.Header.FactoryId, setData.Header.ClassId)
	c.DSet <- setData
	resp := <-msg
	return resp, nil
}

func (c *DataServiceProvider) Update(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	msg := make(chan *protocol.Response, 1)
	defer close(msg)
	setData := SetData{Opt: in.Opt, Data: in.Data, Resp: msg}
	c.DSet <- setData
	resp := <-msg
	return resp, nil
}

func (c *DataServiceProvider) Delete(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	msg := make(chan *protocol.Response, 1)
	defer close(msg)
	setData := SetData{Data: in.Data, Resp: msg}
	c.DSet <- setData
	resp := <-msg
	return resp, nil
}

func (c *DataServiceProvider) Pull(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {

	//stream.Send()
	return nil
}

func (c *DataServiceProvider) Start(dir string) {
	path := fmt.Sprintf("%s/%s", dir, "store")
	core.AppLog.Printf("creating path %s if not existed", path)
	err := os.MkdirAll(path, 0755)
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

	tcp, err := net.Listen("tcp", fmt.Sprintf(":%d", core.RPC_PORT))
	if err != nil {
		panic(err)
	}
	rpc := grpc.NewServer()
	c.server = rpc
	protocol.RegisterDataServiceServer(rpc, c)
	protocol.RegisterPostofficeServiceServer(rpc, c)
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

// internal operations

func (m *DataServiceProvider) create(sd SetData) (KeyIndex, error) {
	ki := sd.IndexKey()
	ki.Header.Revision = 1
	ki.Header.Timestamp = time.Now().UnixMilli()
	ki.Header.Size = int32(len(sd.Value))
	sd.Header.Revision = ki.Header.Revision
	k, v, err := ki.Pair()
	if err != nil {
		return ki, err
	}
	dk, err := sd.DataKey()
	if err != nil {
		return ki, err
	}
	err = m.Local.Db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(k)
		if err == nil {
			return fmt.Errorf("key already existed")
		}
		if err = txn.Set(k, v); err != nil {
			return err
		}
		return txn.Set(dk, sd.Value)
	})
	return ki, err
}
func (m *DataServiceProvider) update(sd SetData) error {
	ki := sd.IndexKey()
	k, err := ki.CompositKey()
	if err != nil {
		return err
	}
	rev := sd.Header.Revision
	return m.Local.Db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return fmt.Errorf("key not existed %s", err.Error())
		}
		if err = item.Value(func(val []byte) error {
			return ki.Val(val)
		}); err != nil {
			return err
		}

		if ki.Header.Revision != rev {
			return fmt.Errorf("revison not matched %d %d", ki.Header.Revision, rev)
		}
		ki.Header.Revision++
		ki.Header.Timestamp = time.Now().UnixMilli()
		ki.Header.Size = int32(len(sd.Value))
		v, err := ki.Value()
		if err != nil {
			return err
		}
		if err = txn.Set(k, v); err != nil {
			return err
		}
		sd.Header.Revision = ki.Header.Revision
		dk, err := sd.DataKey()
		if err != nil {
			return err
		}
		return txn.Set(dk, sd.Value)
	})
}
func (m *DataServiceProvider) delete(sd SetData) error {
	ki := sd.IndexKey()
	k, err := ki.CompositKey()
	if err != nil {
		return err
	}
	return m.Local.Db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return fmt.Errorf("key not existed %s", err.Error())
		}
		if err = item.Value(func(val []byte) error {
			return ki.Val(val)
		}); err != nil {
			return err
		}
		if err = txn.Delete(k); err != nil {
			return err
		}
		sd.Header.Revision = ki.Header.Revision
		dk, err := sd.DataKey()
		if err != nil {
			return err
		}
		return txn.Delete(dk)
	})
}
func (m *DataServiceProvider) get(gd GetData) (*protocol.Data, error) {
	data := protocol.Data{Header: &protocol.Header{}}
	ki := gd.IndexKey()
	k, err := ki.CompositKey()
	if err != nil {
		return &data, err
	}

	err = m.Local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		if err = item.Value(func(val []byte) error {
			return ki.Val(val)
		}); err != nil {
			return nil
		}
		//check revision with requested revision
		gd.Data.Header.Revision = ki.Header.Revision
		dk, err := gd.DataKey()
		if err != nil {
			return err
		}
		item, err = txn.Get(dk)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			data.Value = val
			data.Header.Revision = ki.Header.Revision
			data.Header.Timestamp = ki.Header.Timestamp
			data.Header.Size = ki.Header.Size
			return nil
		})
	})
	return &data, err
}

func (m *DataServiceProvider) reset(sd SetData) error {
	ki := sd.IndexKey()
	k, err := ki.CompositKey()
	if err != nil {
		return err
	}

	return m.Local.Db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return fmt.Errorf("key not existed %s", err.Error())
		}
		if err = item.Value(func(val []byte) error {
			return ki.Val(val)
		}); err != nil {
			return err
		}
		ki.Header.Revision = 1
		ki.Header.Timestamp = time.Now().UnixMilli()
		ki.Header.Size = int32(len(sd.Value))
		v, err := ki.Value()
		if err != nil {
			return err
		}
		if err = txn.Set(k, v); err != nil {
			return err
		}
		sd.Header.Revision = ki.Header.Revision
		dk, err := sd.DataKey()
		if err != nil {
			return err
		}
		return txn.Set(dk, sd.Value)
	})
}
