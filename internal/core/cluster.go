package core

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

type Subscription struct {
	Prefix   string `json:"prefix"`
	Topic    string `json:"topic"`
	Endpoint string `json:"endpoint"`
	Deleting bool   `json:"deleting"`
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

type ClusterService interface {
	HashRing(r RingRequest)
	KeyRing(r RingRequest)
	RingToken(key []byte) uint32
	Request(r DataRequest)
}
