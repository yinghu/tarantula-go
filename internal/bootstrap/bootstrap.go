package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
)

type TarantulaService interface {
	Config() string
	Start(f conf.Env, c cluster.Cluster) error
	Shutdown()
	event.EventService
}

type TarantulaApp interface {
	metrics.MetricsService
	Cluster() cluster.Cluster
	http.Handler

}
