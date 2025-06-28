package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AppAdmin struct {
	TarantulaService
}

func (s *AppAdmin) AccessControl() int32 {
	return ADMIN_ACCESS_CONTROL
}

func (s *AppAdmin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	s.Cluster().OnJoin(core.Node{Name: "a1"})
	w.WriteHeader(http.StatusOK)
	session := core.OnSession{Successful: true, Message: "app admin"}
	w.Write(util.ToJson(session))
}
