package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AppItemAdmin struct {
	TarantulaService
}

func (s *AppItemAdmin) AccessControl() int32 {
	return ADMIN_ACCESS_CONTROL
}

func (s *AppItemAdmin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.WriteHeader(http.StatusOK)
		session := core.OnSession{Successful: true, Message: "app item admin [" + s.Cluster().Group() + "]"}
		w.Write(util.ToJson(session))
		r.Body.Close()
	}()
	cmd := r.PathValue("cmd")
	core.AppLog.Printf("command %s\n", cmd)
	if cmd == "onupdate" {
		s.ItemListener().OnUpdate()
		return
	}
	core.AppLog.Printf("cmd not supported %s\n", cmd)
}
