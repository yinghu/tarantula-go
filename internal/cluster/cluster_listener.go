package cluster

import (
	"sync"

	"gameclustering.com/internal/core"
)

type ClusterListener struct {
	ClusterManager
}

func newListener(group string, etcEndpoints []string, local LocalNode, kl core.ClusterListener) core.Cluster {
	cmanager := ClusterManager{group: group, EtcdEndpoints: etcEndpoints, local: local}
	listener := ClusterListener{ClusterManager: cmanager}
	listener.clistener = kl
	listener.lock = &sync.Mutex{}
	listener.cluster = make(map[string]LocalNode)
	listener.partition = make([]string, core.CLUSTER_PARTITION_NUM)
	listener.quit = make(chan bool)
	return &listener
}

func (c *ClusterListener) Join() error {
	core.AppLog.Printf("Cluster waiting for quit signal %s\n", c.Group)
	<-c.quit
	return nil
}
