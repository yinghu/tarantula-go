package clustering

import (
	context "context"
	"fmt"
	"net"
	"os"
	"sync"

	"gameclustering.com/internal/core"
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
	protocol.UnimplementedTransactionServiceServer
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

	//topic message
	DMessager     chan *protocol.Mail
	subscriptions SubscriptionRegistry
	listeners     map[string]ReceiverAsync //chan *protocol.Topic
	DRequest      chan TopicRequest
}

func (c *DataServiceProvider) Get(ctx context.Context, in *protocol.Request) (*protocol.Response, error) {
	getdata := GetData{in}
	data, err := c.get(getdata)
	if err != nil {
		return &protocol.Response{Successful: false, Message: err.Error()}, err
	}
	return &protocol.Response{Successful: true, Data: &protocol.DataSet{List: []*protocol.Data{data}}}, nil

}

func (c *DataServiceProvider) Query(request *protocol.Request, stream grpc.ServerStreamingServer[protocol.Response]) error {

	tf, existed := core.QueryFactoryRegistry[request.Query.Id]
	if !existed {
		return fmt.Errorf("event factory not registered %s", request.Query.Id)
	}
	q, err := tf().Import(request.Query.Criteria)
	if err != nil {
		return err
	}
	core.AppLog.Debug().Msgf("query %v", q)
	buff := core.NewBuffer(16)
	buff.WriteUInt32(q.QFactoryId())
	buff.WriteUInt32(q.QClassId())
	buff.Flip()
	px, err := buff.Read(0)
	if err != nil {
		return err
	}
	p := px
	dset := make([]*protocol.Data, 0)
	core.AppLog.Debug().Msgf("query : %d %d %d %d", q.QLimit(), q.QOffset(), q.QFactoryId(), q.QClassId())

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
	return stream.Send(&resp)
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
	c.DMessager <- &protocol.Mail{Topic: in, Opt: core.TOPIC_MAIL}
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
	c.DMessager = make(chan *protocol.Mail, NODE_EVENT_BUFFER_SIZE)
	c.DRequest = make(chan TopicRequest, NODE_EVENT_BUFFER_SIZE)
	c.listeners = make(map[string]ReceiverAsync) //chan *protocol.Topic)
	c.subscriptions = SubscriptionRegistry{topicEnds: make(map[core.TopicKey]map[string]core.Subscription), cPools: make(map[core.TopicKey]*core.RpcConnPool)}

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
	protocol.RegisterTransactionServiceServer(rpc, c)
	core.AppLog.Debug().Msgf("local data service provider started on : %s", tcp.Addr().String())
	c.DWait.Done()
	err = rpc.Serve(tcp)
	if err != nil {
		panic(err)
	}
}
