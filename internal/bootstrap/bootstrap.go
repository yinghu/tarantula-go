package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
)

const (
	PUBLIC_ACCESS_CONTROL    int32 = 0
	PROTECTED_ACCESS_CONTROL int32 = 1
	ADMIN_ACCESS_CONTROL     int32 = 6
	SUDO_ACCESS_CONTROL      int32 = 100
)

type TarantulaContext interface {
	Config() string
	Start(f conf.Env, c cluster.Cluster) error
	Shutdown()
	event.EventService
	cluster.KeyListener
}

type TarantulaApp interface {
	Metrics() metrics.MetricsService
	Cluster() cluster.Cluster
	Authenticator() core.Authenticator
	AccessControl() int32
	http.Handler
}
