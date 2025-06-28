package cluster

import (
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
)

func CreateCluster(f conf.Env, kl core.ClusterListener) core.Cluster {
	if f.Clustering {
		return newCluster(f.GroupName, f.EtcdEndpoints, LocalNode{Node: core.Node{Name: f.NodeName, HttpEndpoint: f.HttpEndpoint, TcpEndpoint: f.Evp.TcpEndpoint}}, kl)
	}
	return newListener(f.GroupName, f.EtcdEndpoints, core.Node{Name: f.NodeName, HttpEndpoint: f.HttpEndpoint, TcpEndpoint: f.Evp.TcpEndpoint}, kl)
}
