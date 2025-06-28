package cluster

import (
	"context"
	"slices"
	"sync"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type ClusterListener struct {
	clistener           core.ClusterListener
	Group         string
	EtcdEndpoints []string
	local         core.Node
	lock          *sync.Mutex
	cluster       map[string]core.Node
	partition     []string
	quit          chan bool
}

func newListener(group string, etcEndpoints []string, local core.Node, kl core.ClusterListener) core.Cluster {
	listener := ClusterListener{Group: group, EtcdEndpoints: etcEndpoints, local: local}
	listener.clistener = kl
	listener.lock = &sync.Mutex{}
	listener.cluster = make(map[string]core.Node)
	listener.partition = make([]string, core.CLUSTER_PARTITION_NUM)
	listener.quit = make(chan bool)
	return &listener
}

func (c *ClusterListener) Local() core.Node {
	return c.local
}

func (c *ClusterListener) View() []core.Node {
	c.lock.Lock()
	defer c.lock.Unlock()
	nv := make([]core.Node, 0)
	for _, v := range c.cluster {
		nv = append(nv, v)
	}
	return nv
}

func (c *ClusterListener) Partition(key []byte) core.Node {
	p := util.Partition(key, uint32(core.CLUSTER_PARTITION_NUM))
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cluster[c.partition[p]]
}

func (c *ClusterListener) Atomic(prefix string, t core.Exec) error {
	if prefix == "" {
		prefix = c.Group
		core.AppLog.Printf("Reset Lock prefix %s\n", prefix)
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return err
	}
	defer cli.Close()
	session, err := concurrency.NewSession(cli)
	if err != nil {
		return err
	}
	defer session.Close()
	mutex := concurrency.NewMutex(session, prefix+"#lock")
	ctx := context.Background()
	mutex.Lock(ctx)
	defer mutex.Unlock(ctx)
	return t(&EtcdClient{cli: cli, prefix: prefix})
}

func (c *ClusterListener) OnJoin(join core.Node) {
	c.cluster[join.Name] = join
	c.group()
}

func (c *ClusterListener) Join() error {
	core.AppLog.Printf("Cluster waiting for quit signal %s\n", c.Group)
	<-c.quit
	return nil
}

func (c *ClusterListener) Wait() {

}

func (c *ClusterListener) Quit() {
	c.quit <- true
}

func (c *ClusterListener) Started(){
	
}

func (c *ClusterListener) Listener() core.ClusterListener {
	return c.clistener
}
func (c *ClusterListener) group() {
	sz := len(c.cluster)
	core.AppLog.Printf("Cluster grouping %d\n", sz)
	nds := make([]string, sz)
	i := 0
	for n := range c.cluster {
		nds[i] = n
		i++
	}
	slices.Sort(nds)
	for p := range core.CLUSTER_PARTITION_NUM {
		i := p % sz
		c.partition[p] = nds[i]
	}
}
