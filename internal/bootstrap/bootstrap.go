package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/metrics"
)

const (
	PUBLIC_ACCESS_CONTROL    int32 = 0
	PROTECTED_ACCESS_CONTROL int32 = 1
	ADMIN_ACCESS_CONTROL     int32 = 30
	SUDO_ACCESS_CONTROL      int32 = 100
)

const (
	DB_OP_ERR_CODE int = 500100

	WRONG_PASS_CODE int    = 400100
	WRONG_PASS_MSG  string = "wrong user/password"

	ILLEGAL_TOKEN_CODE  int = 400101
	INVALID_TOKEN_CODE  int = 400102
	ILLEGAL_ACCESS_CODE int = 400103

	INVALID_TOKEN_MSG string = "invalid token"
)


type TarantulaContext interface {
	Config() string
	Start(f conf.Env, c cluster.Cluster) error
	Shutdown()
	event.EventService
	cluster.KeyListener
}

type TarantulaApp interface {
	ItemService() item.ItemService
	Metrics() metrics.MetricsService
	Cluster() cluster.Cluster
	Authenticator() core.Authenticator
	AccessControl() int32
	Request(sesion core.OnSession, w http.ResponseWriter, r *http.Request)
}
