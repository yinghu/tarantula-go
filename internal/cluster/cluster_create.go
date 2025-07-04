package cluster

import (
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

func CreateCluster(f conf.Env, kl core.ClusterListener) core.Cluster {
	node := core.Node{Name: f.NodeName, HttpEndpoint: f.HttpEndpoint, TcpEndpoint: f.Evp.TcpEndpoint}
	if f.Clustering {
		return newCluster(f.ClusterCtx(), f.EtcdEndpoints, LocalNode{Node: node}, kl)
	}
	return newListener(f.ClusterCtx(), f.EtcdEndpoints, LocalNode{Node: node}, kl)
}
