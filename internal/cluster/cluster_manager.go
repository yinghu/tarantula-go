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

type ClusterManager struct {
	clistener     core.ClusterListener
	group         string
	EtcdEndpoints []string
	local         LocalNode
	lock          *sync.Mutex
	cluster       map[string]LocalNode
	partition     []string
	quit          chan bool
}

func (c *ClusterManager) Local() core.Node {
	return c.local.Node
}

func (c *ClusterManager) View() []core.Node {
	c.lock.Lock()
	defer c.lock.Unlock()
	nv := make([]core.Node, 0)
	for _, v := range c.cluster {
		nv = append(nv, v.Node)
	}
	return nv
}

func (c *ClusterManager) Partition(key []byte) core.Node {
	p := util.Partition(key, uint32(core.CLUSTER_PARTITION_NUM))
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.cluster[c.partition[p]].Node
}

func (c *ClusterManager) Atomic(prefix string, t core.Exec) error {
	if prefix == "" {
		prefix = c.Group()
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

func (c *ClusterManager) Listener() core.ClusterListener {
	return c.clistener
}

func (c *ClusterManager) Quit() {
	c.quit <- true
}

func (c *ClusterManager) grouping() {
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

func (c *ClusterManager) Wait() {

}

func (c *ClusterManager) Started() {

}

func (c *ClusterManager) OnJoin(join core.Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, joined := c.cluster[join.Name]
	if joined {
		return
	}
	c.cluster[join.Name] = LocalNode{Node: join}
	c.grouping()
}

func (c *ClusterManager) OnLeave(leave core.Node) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.cluster, leave.Name)
	c.grouping()
}

func (c *ClusterManager) OnUpdate(key string, value string, opt core.Opt) {
	core.AppLog.Printf("Key update %s, %s, %v\n", key, value, opt)
}

func (c *ClusterManager) Group() string {
	return c.group
}
