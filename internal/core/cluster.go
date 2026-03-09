package core

const (
	RPC_PORT int = 7001

	SET_OPT_RECOVER int32 = 1
	SET_OPT_CLOSE   int32 = 2

	GET_DATA_REQUEST    int32 = 10
	CREATE_DATA_REQUEST int32 = 11
	UPDATE_DATA_REQUEST int32 = 12
	DELETE_DATA_REQUEST int32 = 13
	RESET_DATA_REQUEST  int32 = 14
	
)

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

type RingSync struct {
	Remote string   `json:"remote"`
	Hashs  []uint32 `json:"hashs"`
}

type RingRequest struct {
	Opt      int32
	Address  string
	Source   RingSync
	Token    uint32
	Replicas int
	Async    chan []Node
}

type DataHeader struct {
	FactoryId int32
	ClassId   int32
	Revision  int64
}


type DataRequest struct {
	DataHeader
	Prefix int32
	Key    []byte
	Value  []byte
	Opt    int32
	Async  chan Chunk
}

type ClusterService interface {
	HashRing(r RingRequest)
	KeyRing(r RingRequest)
	RingToken(key []byte) uint32
	Request(r DataRequest)
}
