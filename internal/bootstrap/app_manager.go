package bootstrap


import (
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/metrics"
)

type AppManager struct{
	Cls     cluster.Cluster
	Metr    metrics.MetricsService
	Auth    core.Authenticator
}

func (s *AppManager) Metrics() metrics.MetricsService {
	return s.Metr
}
func (s *AppManager) Cluster() cluster.Cluster {
	return s.Cls
}
func (s *AppManager) Authenticator() core.Authenticator {
	return s.Auth
}