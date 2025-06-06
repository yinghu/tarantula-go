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

const (
	DB_OP_ERR_CODE     int    = 500100
	WRONG_PASS_CODE    int    = 400100
	WRONG_PASS_MSG     string = "wrong user/password"
	INVALID_TOKEN_CODE int    = 400101
	INVALID_TOKEN_MSG  string = "invalid token"
)

const (
	METRICS_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS req_metrics (id BIGSERIAL PRIMARY KEY,path VARCHAR(50) NOT NULL,req_timed BIGINT NOT NULL,req_time TIMESTAMP DEFAULT NOW(),node VARCHAR(10) NOT NULL,req_id INTEGER DEFAULT 0)"
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
	Request(sesion core.OnSession, w http.ResponseWriter, r *http.Request)
}
