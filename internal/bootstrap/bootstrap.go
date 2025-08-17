package bootstrap

import (
	"net/http"

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
	STANDALONE_APP     int    = 200000
	STANDALONE_APP_MSG string = "publish skipped"
	DB_OP_ERR_CODE     int    = 500100

	WRONG_PASS_CODE int    = 400100
	WRONG_PASS_MSG  string = "wrong user/password"

	BAD_REQUEST_CODE    int = 400100
	ILLEGAL_TOKEN_CODE  int = 400101
	INVALID_TOKEN_CODE  int = 400102
	ILLEGAL_ACCESS_CODE int = 400103
	INVALID_TICKET_CODE int = 400104
	INVALID_JSON_CODE   int = 400105

	INVALID_TOKEN_MSG  string = "invalid token"
	ILLEGAL_ACCESS_MSG string = "illegal access"
	ILLEGAL_TOKEN_MSG  string = "bad token"
	BAD_REQUEST_MSG    string = "bad request"
)

type TarantulaContext interface {
	Config() string
	Start(f conf.Env, c core.Cluster) error
	Shutdown()
	event.EventService
	core.ClusterListener
	Context() string
	Service() TarantulaService
}

type TarantulaService interface {
	ItemService() item.ItemService
	Metrics() metrics.MetricsService
	Cluster() core.Cluster
	Authenticator() core.Authenticator
	Sequence() core.Sequence
	ItemListener() item.ItemListener
	BootstrapListener() BootstrapListener
}

type TarantulaApp interface {
	TarantulaService
	AccessControl() int32
	Request(sesion core.OnSession, w http.ResponseWriter, r *http.Request)
}

type BootstrapListener interface {
	NodeStarted(n core.Node)
}

type Login struct {
	Id            int32            `json:"-"`
	Name          string           `json:"login"`
	Hash          string           `json:"password"`
	ReferenceId   int32            `json:"referenceId"`
	SystemId      int64            `json:"systemId:string"`
	AccessControl int32            `json:"accessControl,string"`
	Cc            chan event.Chunk `json:"-"`
}
