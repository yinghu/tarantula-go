package core

const (
	CLUSTER_PARTITION_NUM int = 271
)

type Node struct {
	Name         string `json:"name"`
	RingToken    uint32 `json:"ringToken"`
	Meta         string `json:"meta"`
	IP           string `json:"address"`
	State        int    `json:"state"`
	HttpEndpoint string `json:"http"`
	TcpEndpoint  string `json:"tcp"`
}
type KVLoad func(k, v string) bool

type Ctx interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Del(key string, withPrefix bool) error
	List(prefix string, loaded KVLoad) error
}

type Exec func(ctx Ctx) error

type Cluster interface {
	Group() string
	Local() Node
	View() []Node
	Partition(key []byte) Node
	Atomic(prefix string, t Exec) error
	Join() error
	Wait()
	Quit()
	Started()

	OnJoin(join Node)
	OnLeave(leave Node)
	Listener() ClusterListener
}

type Opt struct {
	IsCreate bool   `json:"IsCreate"`
	IsModify bool   `json:"IsModify"`
	Type     string `json:"Type"`
}

type ClusterListener interface {
	KVUpdated(key string, value string, opt Opt)
	MemberJoined(joined Node)
	MemberLeft(left Node)
}

type RingRequest struct {
	Token    uint32
	Replicas int
	Async    chan []Node
}

type GetRequest struct {
	Database string
	Key      []byte
	Async    chan Chunk
}
type SetRequest struct {
	Database string
	Key      []byte
	Value    []byte
	Async    chan Chunk
}

type ClusterService interface {
	Sequence
	HashRing(r RingRequest)
	KeyRing(r RingRequest)
	RingToken(key []byte) uint32
	Get(get GetRequest)
	Set(get SetRequest)
}
