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
	Local() Node
	View() []Node
	Partition(key []byte) Node
	Atomic(prefix string, t Exec) error
	Join() error
	Wait()
	Quit()

	OnJoin(join Node)
	Listener() KeyListener
}

type KeyListener interface {
	Updated(key string, value string)
}
