package core

const (
	CLUSTER_PARTITION_NUM int = 271
)

type Node struct {
	Name         string `json:"name"`
	HttpEndpoint string `json:"http"`
	TcpEndpoint  string `json:"tcp"`
}

type Ctx interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Del(key string) error
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

type ClusterListener interface {
	Updated(key string, value string)
	MemberJoined(joined Node)
	MemberLeft(left Node)
}
