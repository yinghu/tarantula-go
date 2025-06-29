package bootstrap

import (
	"encoding/json"
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
	cmd := r.PathValue("cmd")
	core.AppLog.Printf("Call on cmd %s %v\n", cmd,s.Cluster().Local())
	if cmd == "join" {
		var join core.Node
		json.NewDecoder(r.Body).Decode(&join)
		s.Cluster().OnJoin(join)
		core.AppLog.Printf("Call on join %s\n", join.Name)
	}
	w.WriteHeader(http.StatusOK)
	session := core.OnSession{Successful: true, Message: "app admin"}
	w.Write(util.ToJson(session))
}
