package main

import (
	context "context"
	"fmt"
	"net"
	"os"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/tarantula/protocol"
	badger "github.com/dgraph-io/badger/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type DataServiceProvider struct {
	protocol.UnimplementedDataServiceServer
	Local  *persistence.BadgerLocal
	RNode  <-chan []core.Node
	server *grpc.Server
}

func (c *DataServiceProvider) Get(ctx context.Context, in *protocol.Request) (*protocol.Data, error) {
	data := protocol.Data{}
	err := c.Local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(in.Key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			data.Value = val
			return nil
		})
	})
	return &data, err
}

func (c *DataServiceProvider) Set(ctx context.Context, in *protocol.Data) (*protocol.Response, error) {
	err := c.Local.Db.Update(func(txn *badger.Txn) error {
		return txn.Set(in.Key, in.Value)
	})
	if err != nil {
		return &protocol.Response{Successful: false, Message: err.Error()}, err
	}
	return &protocol.Response{Successful: true, Message: "saved"}, nil
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
	tcp, err := net.Listen("tcp", ":7001")
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

func (m *DataServiceProvider) ClientGet(target *core.Node, request *core.GetRequest) ([]byte, error) {
	tcp, err := grpc.NewClient(target.IP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	var dt *protocol.Data
	dt, err = dsp.Get(context.Background(), &protocol.Request{Key: request.Key})
	if err != nil {
		return nil, err
	}
	return dt.Value, nil
}

func (m *DataServiceProvider) ClientSet(target *core.Node, request *core.SetRequest) (*protocol.Response, error) {
	//core.AppLog.Printf("target node %s %s %d\n", ringNode.IP, ringNode.Name, ringNode.RingToken)
	tcp, err := grpc.NewClient(target.IP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &protocol.Response{}, err
	}
	defer tcp.Close()
	dsp := protocol.NewDataServiceClient(tcp)
	kv := protocol.Data{Key: request.Key, Value: request.Value}
	return dsp.Set(context.Background(), &kv)
}
func (m *DataServiceProvider) RingUpdated() {
	stopping := false
	for nlist := range m.RNode {
		for _, n := range nlist {
			if n.State == -1000 {
				stopping = true
				break
			}
			core.AppLog.Printf("node updated IP : %s NAME : %s RING TOKEN : %d STATE : %d", n.IP, n.Name, n.RingToken, n.State)
		}
		if stopping {
			break
		}
	}
	m.server.Stop()
	m.Local.Close()
	core.AppLog.Info().Msg("local data service provider has stopped")
}

// internal operations
func (m *DataServiceProvider) SaveKeyIndex(keyIndex *KeyIndex) error {
	kBuff := core.NewBuffer(100)
	if err := keyIndex.WriteKey(kBuff); err != nil {
		return err
	}
	kBuff.Flip()
	key, err := kBuff.Read(0)
	if err != nil {
		return err
	}
	kBuff.Clear()
	if err := keyIndex.Write(kBuff); err != nil {
		return err
	}
	kBuff.Flip()
	value, err := kBuff.Read(0)
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

func (m *DataServiceProvider) LoadKeyIndex(keyIndex *KeyIndex) error {
	kBuff := core.NewBuffer(100)
	if err := keyIndex.WriteKey(kBuff); err != nil {
		return err
	}
	kBuff.Flip()
	key, err := kBuff.Read(0)
	if err != nil {
		return err
	}
	kBuff.Clear()
	if err := m.Local.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			if err := kBuff.Write(val); err != nil {
				return err
			}
			return nil
		})
	}); err != nil {
		return err
	}
	kBuff.Flip()
	return keyIndex.Read(kBuff)
}
