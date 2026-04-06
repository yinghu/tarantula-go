package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/protocol"
	"github.com/rs/zerolog"
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

type LogForwarder interface {
	Forward(level zerolog.Level, log []byte)
}

type TarantulaContext interface {
	Config() string
	Start(f core.Env) error
	Shutdown()
	Context() string
	Service() TarantulaService
}

type TarantulaService interface {
	ItemService() item.ItemService
	Metrics() core.MetricsService
	Authenticator() core.Authenticator
	Sequence() core.Sequence
	ItemListener() item.ItemListener
	Pusher() core.Pusher
	Cluster() core.ClusterService
	Event() core.EventService
	RegisterLogForwarder(logf LogForwarder)
}

type TarantulaApp interface {
	TarantulaService
	AccessControl() int32
	NodeId() string
	Context() string
	Forward(topic *protocol.Topic)
	Request(sesion core.OnSession, w http.ResponseWriter, r *http.Request)
}

type Login struct {
	Id            int32           `json:"-"`
	Name          string          `json:"login"`
	Hash          string          `json:"password"`
	ReferenceId   int32           `json:"referenceId"`
	SystemId      int64           `json:"systemId:string"`
	AccessControl int32           `json:"accessControl,string"`
	Cc            chan core.Chunk `json:"-"`
}
