package bootstrap

import (
	"encoding/json"
	"net/http"
	"strings"

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
	defer func() {
		w.WriteHeader(http.StatusOK)
		session := core.OnSession{Successful: true, Message: "app admin"}
		w.Write(util.ToJson(session))
		r.Body.Close()
	}()
	cmd := r.PathValue("cmd")
	core.AppLog.Printf("Call on cmd %s %v\n", cmd, s.Cluster().Local())
	if cmd == "join" {
		var join core.Node
		json.NewDecoder(r.Body).Decode(&join)
		s.Cluster().OnJoin(s.convert(join))
		core.AppLog.Printf("Call on join %v\n", join)
		return
	}
	if cmd == "left" {
		var left core.Node
		json.NewDecoder(r.Body).Decode(&left)
		s.Cluster().OnLeave(s.convert(left))
		core.AppLog.Printf("Call on left %v\n", left)
	}
}

func (s *AppAdmin) convert(node core.Node) core.Node {
	node.Name = strings.Replace(node.Name, "admin", s.Cluster().Group(), 1)
	lparts := strings.Split(s.Cluster().Local().TcpEndpoint, ":")
	rparts := strings.Split(node.TcpEndpoint, ":")
	node.TcpEndpoint = rparts[0] + ":" + rparts[1] + ":" + lparts[2]
	core.AppLog.Printf("Node : %v\n", node)
	return node
}
