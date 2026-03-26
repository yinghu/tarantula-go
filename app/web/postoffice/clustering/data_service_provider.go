package clustering

import (
	context "context"
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

const (
	PULL_BATCH_SIZE int = 10
)

type DataServiceProvider struct {
	protocol.UnimplementedDataServiceServer
	protocol.UnimplementedPostofficeServiceServer
	Local       *persistence.BadgerLocal
	RNode       <-chan RingUpdate
	RSync       <-chan []byte
	server      *grpc.Server
	Mll         *MemberListListener
	backRing    NodeRing
	rpcEndpoint string
	//write worker chan
	DSet    chan SetData
	DPull   chan core.RingSync
	DWait   sync.WaitGroup
	running bool

	//messaging
	DMessager chan *protocol.Topic
	//index         map[string]core.Subscription
	subscriptions SubscriptionRegistry
	listeners     map[string]chan *protocol.Topic
	DRequest      chan TopicRequest
}

func (c *DataServiceProvider) Get(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {
	if request.Opt == core.GET_DATA_REQUEST {
		getdata := GetData{request}
		data, err := c.get(getdata)
		if err != nil {
			return err
		}
		resp := protocol.Response{Successful: true, Data: &protocol.DataSet{List: []*protocol.Data{data}}}
		return stream.Send(&resp)
	}
	if request.Opt == core.QUERY_DATA_REQUEST {
		q := event.CreateQuery(request.Query.Id)
		err := event.Import(q, request.Query.Criteria, 100)
		if err != nil {
			return err
		}
		buff := core.NewBuffer(16)
		buff.WriteUInt32(q.QFactoryId())
		buff.WriteUInt32(q.QClassId())
		buff.Flip()
		px, err := buff.Read(0)
		if err != nil {
			return err
		}
		p := px
		rc := make(chan *protocol.Response, 3)
		dset := make([]*protocol.Data, 0)
		core.AppLog.Debug().Msgf("query : %d %d ", q.QLimit(), q.QOffset())
		go func() {
			limit := q.QLimit()
			offset := q.QOffset()
			c.Local.Db.View(func(txn *badger.Txn) error {
				op := badger.IteratorOptions{PrefetchSize: 100, PrefetchValues: false, Reverse: false}
				it := txn.NewIterator(op)
				defer it.Close()
				for it.Seek(p); it.ValidForPrefix(p); it.Next() {
					if offset > 0 {
						offset--
						continue
					}
					p = px
					item := it.Item()
					k := append([]byte{}, item.Key()[12:]...)
					item.Value(func(val []byte) error {
						if q.QFilter(k, val) {
							v := append([]byte{}, val...)
							dset = append(dset, &protocol.Data{Key: k, Value: v, Header: &protocol.Header{}})
							limit--
						}
						return nil
					})
					if limit == 0 {
						break
					}
				}
				return nil
			})
			resp := protocol.Response{Successful: true, Data: &protocol.DataSet{List: dset}}
			rc <- &resp
		}()
		rs := <-rc
		return stream.Send(rs)
	}
	return fmt.Errorf("opt not supported %d", request.Opt)
}

func (c *DataServiceProvider) Reset(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	msg := make(chan *protocol.Response, 1)
	defer close(msg)
	setData := SetData{Opt: in.Opt, Data: in.Data, Resp: msg}
	c.DSet <- setData
	resp := <-msg
	return resp, nil
}

func (c *DataServiceProvider) Create(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	msg := make(chan *protocol.Response, 1)
	defer close(msg)
	setData := SetData{Opt: in.Opt, Prefix: in.Prefix, Data: in.Data, Resp: msg}
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
	setData := SetData{Opt: in.Opt, Data: in.Data, Resp: msg}
	c.DSet <- setData
	resp := <-msg
	return resp, nil
}

func (c *DataServiceProvider) Pull(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {
	ch := make(chan *protocol.Response, 3)
	go c.pull(request.Prefix, request.Opt, ch)
	for resp := range ch {
		if !resp.Successful {
			break
		}
		if err := stream.Send(resp); err != nil {
			break
		}
	}
	return nil
}

func (c *DataServiceProvider) Send(ctx context.Context, in *protocol.Topic) (*protocol.Response, error) {
	c.DMessager <- in
	return &protocol.Response{Successful: true, Message: "event published"}, nil
}

func (c *DataServiceProvider) Start(dir string) {
	c.running = true
	c.backRing = NodeRing{nodes: make([]core.Node, 0)}
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
	c.DMessager = make(chan *protocol.Topic, NODE_EVENT_BUFFER_SIZE)
	c.DRequest = make(chan TopicRequest, NODE_EVENT_BUFFER_SIZE)
	c.listeners = make(map[string]chan *protocol.Topic)
	c.subscriptions = SubscriptionRegistry{Ends: make(map[core.TopicKey]map[string]core.Subscription)}

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
	c.DWait.Done()
	err = rpc.Serve(tcp)
	if err != nil {
		panic(err)
	}
}

// internal operations

func (m *DataServiceProvider) create(sd SetData) (KeyIndex, error) {
	ki := sd.IndexKey()
	ki.Header.Revision = 1
	ki.Header.Timestamp = uint64(time.Now().UnixMilli())
	ki.Header.Size = uint32(len(sd.Value))
	ki.Header.Mutable = sd.Header.Mutable
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
		item, err := txn.Get(k)
		if err == nil {
			if item.Value(func(val []byte) error {
				ev := append([]byte{}, val...)
				return ki.Val(ev)
			}) != nil || ki.Header.State != core.DATA_STATE_DELETED {
				return fmt.Errorf("key already existed")
			}
		}
		if err = txn.Set(k, v); err != nil {
			return err
		}
		return txn.Set(dk, sd.Value)
	})
	return ki, err
}
func (m *DataServiceProvider) update(sd SetData) (KeyIndex, error) {
	ki := sd.IndexKey()
	k, err := ki.CompositKey()
	if err != nil {
		return ki, err
	}
	rev := sd.Header.Revision
	err = m.Local.Db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return fmt.Errorf("key not existed %s", err.Error())
		}
		if err = item.Value(func(val []byte) error {
			v := append([]byte{}, val...)
			return ki.Val(v)
		}); err != nil {
			return err
		}
		if !ki.Header.Mutable {
			return fmt.Errorf("cannot update on immutable %v", ki.Header.Mutable)
		}
		if ki.Header.Revision != rev {
			return fmt.Errorf("revison not matched %d %d", ki.Header.Revision, rev)
		}
		ki.Header.Revision++
		ki.Header.Timestamp = uint64(time.Now().UnixMilli())
		ki.Header.Size = uint32(len(sd.Value))
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
	return ki, err
}
func (m *DataServiceProvider) delete(sd SetData) (KeyIndex, error) {
	ki := sd.IndexKey()
	k, err := ki.CompositKey()
	if err != nil {
		return ki, err
	}
	err = m.Local.Db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return fmt.Errorf("key not existed %s", err.Error())
		}
		if err = item.Value(func(val []byte) error {
			v := append([]byte{}, val...)
			return ki.Val(v)
		}); err != nil {
			return err
		}
		if !ki.Header.Mutable {
			return fmt.Errorf("cannot delete on immutable %v", ki.Header.Mutable)
		}
		if ki.Header.State == core.DATA_STATE_DELETED {
			return fmt.Errorf("already deleted on %d", ki.Header.State)
		}
		ki.Header.State = core.DATA_STATE_DELETED
		ki.Header.Revision++
		uv, err := ki.Value()
		if err != nil {
			return err
		}
		if err = txn.Set(k, uv); err != nil {
			return err
		}
		sd.Header.Revision = ki.Header.Revision
		dk, err := sd.DataKey()
		if err != nil {
			return err
		}
		return txn.Delete(dk)
	})
	return ki, err
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
			v := append([]byte{}, val...)
			return ki.Val(v)
		}); err != nil {
			return err
		}
		//use latest revision
		if ki.Header.State == core.DATA_STATE_DELETED {
			return fmt.Errorf("data deleted %d", ki.Header.State)
		}
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
			v := append([]byte{}, val...)
			data.Value = v
			data.Header.Revision = ki.Header.Revision
			data.Header.Timestamp = ki.Header.Timestamp
			data.Header.Size = ki.Header.Size
			data.Header.Mutable = ki.Header.Mutable
			return nil
		})
	})
	return &data, err
}

func (m *DataServiceProvider) reset(sd SetData) (KeyIndex, error) {
	ki := sd.IndexKey()
	k, err := ki.CompositKey()
	if err != nil {
		return ki, err
	}
	err = m.Local.Db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(k)
		if err != nil {
			return fmt.Errorf("key not existed %s", err.Error())
		}
		if err = item.Value(func(val []byte) error {
			v := append([]byte{}, val...)
			return ki.Val(v)
		}); err != nil {
			return err
		}
		if !ki.Header.Mutable {
			return fmt.Errorf("cannot reset on immutable %v", ki.Header.Mutable)
		}
		ki.Header.Revision = 1
		ki.Header.Timestamp = uint64(time.Now().UnixMilli())
		ki.Header.Size = uint32(len(sd.Value))
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
	return ki, err
}

func (m *DataServiceProvider) pull(from, to uint32, ch chan *protocol.Response) {
	index := KeyIndex{}
	pre, _ := index.lookupPrefix(INDEX_PREFIX)
	data := make([]*protocol.Data, 0, PULL_BATCH_SIZE)
	m.Local.Db.View(func(txn *badger.Txn) error {
		op := badger.IteratorOptions{PrefetchSize: 100, PrefetchValues: false, Reverse: false}
		it := txn.NewIterator(op)
		defer it.Close()
		for it.Seek(pre); it.ValidForPrefix(pre); it.Next() {
			item := it.Item()
			k := append([]byte{}, item.Key()...)
			var v []byte
			item.Value(func(val []byte) error {
				v = append(v, val...)
				return nil
			})
			ki := KeyIndex{Header: &protocol.Header{}}
			err := core.Import(&ki, k, v, 300)
			if err != nil {
				core.AppLog.Warn().Msgf("should not be a error from import data %s", err.Error())
				continue
			}
			if ki.Header.State == core.DATA_STATE_DELETED {
				key, _ := ki.lookupDataKey()
				kz := len(key)
				vdata := protocol.Data{Key: key[12 : kz-8], Header: &protocol.Header{FactoryId: ki.Header.FactoryId, ClassId: ki.Header.ClassId, Revision: ki.Header.Revision, Timestamp: ki.Header.Timestamp, Mutable: ki.Header.Mutable, State: ki.Header.State}}
				data = append(data, &vdata)
				if len(data) == PULL_BATCH_SIZE {
					resp := protocol.Response{Successful: true, Data: &protocol.DataSet{List: data}}
					ch <- &resp
					data = make([]*protocol.Data, 0, PULL_BATCH_SIZE)
				}
				continue
			}
			if from < to {
				if ki.Prefix >= from && ki.Prefix < to {
					key, _ := ki.lookupDataKey()
					vitem, err := txn.Get(key)
					if err != nil {
						core.AppLog.Warn().Msgf("should not be a value missing %v", ki)
						continue
					}
					vitem.Value(func(val []byte) error {
						fv := append([]byte{}, val...)
						kz := len(key)
						vdata := protocol.Data{Key: key[12 : kz-8], Value: fv, Header: &protocol.Header{FactoryId: ki.Header.FactoryId, ClassId: ki.Header.ClassId, Revision: ki.Header.Revision, Timestamp: ki.Header.Timestamp, Mutable: ki.Header.Mutable, State: ki.Header.State}}

						data = append(data, &vdata)
						if len(data) == PULL_BATCH_SIZE {
							resp := protocol.Response{Successful: true, Data: &protocol.DataSet{List: data}}
							ch <- &resp
							data = make([]*protocol.Data, 0, PULL_BATCH_SIZE)
						}
						return nil
					})
				}
			} else {

				if ki.Prefix >= from || ki.Prefix < to {
					key, _ := ki.lookupDataKey()
					vitem, err := txn.Get(key)
					if err != nil {
						core.AppLog.Warn().Msgf("should not be a value missing %v", ki)
						continue
					}
					vitem.Value(func(val []byte) error {
						fv := append([]byte{}, val...)
						kz := len(key)
						vdata := protocol.Data{Key: key[12 : kz-8], Value: fv, Header: &protocol.Header{FactoryId: ki.Header.FactoryId, ClassId: ki.Header.ClassId, Revision: ki.Header.Revision, Timestamp: ki.Header.Timestamp, Mutable: ki.Header.Mutable, State: ki.Header.State}}
						data = append(data, &vdata)
						if len(data) == PULL_BATCH_SIZE {
							resp := protocol.Response{Successful: true, Data: &protocol.DataSet{List: data}}
							ch <- &resp
							data = make([]*protocol.Data, 0, PULL_BATCH_SIZE)
						}
						return nil
					})
				}
			}
		}
		return nil
	})
	if len(data) > 0 {

		resp := protocol.Response{Successful: true, Data: &protocol.DataSet{List: data}}
		ch <- &resp
	}
	ch <- &protocol.Response{Successful: false}
}

func (c *DataServiceProvider) set(resp *protocol.Response) {
	err := c.Local.Db.Update(func(txn *badger.Txn) error {
		for _, d := range resp.Data.List {
			setdata := SetData{Data: d}
			ki := setdata.IndexKey()
			k, v, err := ki.Pair()
			if err != nil {
				core.AppLog.Warn().Msgf("wrong key index %s", err.Error())
				continue
			}
			dkey, err := setdata.DataKey()
			if err != nil {
				core.AppLog.Warn().Msgf("wrong data key %s", err.Error())
				continue
			}
			item, err := txn.Get(k)
			if err != nil { //no data
				txn.Set(k, v)
				if ki.Header.State != core.DATA_STATE_DELETED {
					txn.Set(dkey, setdata.Value)
				}
				continue
			}
			item.Value(func(val []byte) error {
				eki := KeyIndex{Header: &protocol.Header{}}
				ev := append([]byte{}, val...)
				err = core.Import(&eki, k, ev, COMPOSIT_KEY_MAX)
				if err != nil {
					return err
				}
				if eki.Header.Revision < ki.Header.Revision {
					txn.Set(k, v)
					if ki.Header.State == core.DATA_STATE_DELETED {
						txn.Delete(dkey)
					} else {
						txn.Set(dkey, setdata.Value)
					}
				}
				return nil
			})
		}
		return nil
	})
	if err != nil {
		core.AppLog.Warn().Msgf("set err %s", err.Error())
		return
	}
}
