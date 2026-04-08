package core

import (
	"fmt"

	"gameclustering.com/internal/protocol"
	"google.golang.org/grpc"
)

const (
	RPC_PORT int = 7001

	SET_OPT_RECOVER uint32 = 1
	SET_OPT_CLOSE   uint32 = 2

	GET_DATA_REQUEST    uint32 = 10
	CREATE_DATA_REQUEST uint32 = 11
	UPDATE_DATA_REQUEST uint32 = 12
	DELETE_DATA_REQUEST uint32 = 13
	RESET_DATA_REQUEST  uint32 = 14
	QUERY_DATA_REQUEST  uint32 = 15

	PULL_DATA_REQUEST uint32 = 16

	DATA_STATE_READY   uint32 = 0
	DATA_STATE_PENDING uint32 = 1
	DATA_STATE_DELETED uint32 = 2

	COMPOSIT_KEY_MAX int = 500
)

type Chunk struct {
	Remaining bool
	Data      any
}
type Node struct {
	Name         string `json:"name"`
	RingToken    uint32 `json:"ringToken"`
	Meta         string `json:"meta"`
	IP           string `json:"address"`
	State        int    `json:"-"`
	RpcEndpoint  string `json:"rpc,omitempty"`
	HttpEndpoint string `json:"http,omitempty"`
	TcpEndpoint  string `json:"tcp,omitempty"`

	CPool *RpcConnPool `json:"-"`
}
type KVLoad func(k, v string) bool

type Ctx interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Del(key string, withPrefix bool) error
	List(prefix string, loaded KVLoad) error
}

type Exec func(ctx Ctx) error

type Opt struct {
	IsCreate bool   `json:"IsCreate"`
	IsModify bool   `json:"IsModify"`
	Type     string `json:"Type"`
}

type RingRange struct {
	//range >= from and < to
	From uint32 `json:"from"`
	To   uint32 `json:"to"`
}

type TopicKey struct {
	Topic    string
	Endpoint string
}

type Subscription struct {
	NodeId   string `json:"nodeId"`
	Tag      string `json:"tag"`
	Topic    string `json:"topic"`
	Endpoint string `json:"endpoint"`
	Deleting bool   `json:"deleting"`

	CPool *RpcConnPool `json:"-"`
}

func (s *Subscription) Key() string {
	return fmt.Sprintf("%s:%s:%s", s.NodeId, s.Tag, s.Topic)
}
func (s *Subscription) TopicKey() TopicKey {
	return TopicKey{Topic: s.Topic, Endpoint: s.Endpoint}
}

type RingSync struct {
	Remote string       `json:"remote"`
	Ranges []RingRange  `json:"ranges"`
	Sub    Subscription `json:"sub"`
}

type RingRequest struct {
	Opt      uint32
	Address  string
	Source   RingSync
	Token    uint32
	Replicas int
	Async    chan []Node
}

type DataHeader struct {
	FactoryId uint32
	ClassId   uint32
	Revision  uint64
	Mutable   bool
	State     uint32
}

type DataRequest struct {
	DataHeader
	Prefix   uint32
	Key      []byte
	Value    []byte
	Opt      uint32
	Criteria Query
	Async    chan Chunk
}

type TopicListener interface {
	OnTopic(topic *protocol.Topic) error
}

type ClusterService interface {
	HashRing(r RingRequest) (grpc.ServerStreamingClient[protocol.HashNode], error)
	KeyRing(r RingRequest) (grpc.ServerStreamingClient[protocol.HashNode], error)
	RingToken(key []byte) uint32
	Request(r DataRequest)

	Publish(e *protocol.Topic) error
	//List(q Query)
	Subscribe(topic string, listener TopicListener) error
	Unsubscribe(topic string) error
}
