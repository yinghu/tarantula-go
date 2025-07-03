package bootstrap

import (
	"encoding/json"
	"net/http"
	"strings"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AppClusterAdmin struct {
	TarantulaService
}

func (s *AppClusterAdmin) AccessControl() int32 {
	return ADMIN_ACCESS_CONTROL
}

func (s *AppClusterAdmin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cmd := r.PathValue("cmd")
	core.AppLog.Printf("Call on cmd %s %v\n", cmd, s.Cluster().Local())
	session := core.OnSession{Successful: true, Message: "app cluster admin [" + s.Cluster().Group() + "]"}
	if cmd == "join" {
		var join core.Node
		json.NewDecoder(r.Body).Decode(&join)
		s.Cluster().OnJoin(s.convert(join))
		core.AppLog.Printf("Call on join %v\n", join)
		w.WriteHeader(http.StatusOK)
		w.Write(util.ToJson(session))
		return
	}
	if cmd == "left" {
		var left core.Node
		json.NewDecoder(r.Body).Decode(&left)
		s.Cluster().OnLeave(s.convert(left))
		core.AppLog.Printf("Call on left %v\n", left)
		w.WriteHeader(http.StatusOK)
		w.Write(util.ToJson(session))
		return
	}
	if cmd == "update" {
		var update KVUpdate
		json.NewDecoder(r.Body).Decode(&update)
		core.AppLog.Printf("%s, %s, %s\n", update.Key, update.Value, update.Type)
		w.WriteHeader(http.StatusOK)
		w.Write(util.ToJson(session))
		return
	}
	core.AppLog.Printf("cmd not supported %s\n", cmd)
	w.WriteHeader(http.StatusOK)
	w.Write(util.ToJson(session))
}

func (s *AppClusterAdmin) convert(node core.Node) core.Node {
	node.Name = strings.Replace(node.Name, "admin", s.Cluster().Group(), 1)
	lparts := strings.Split(s.Cluster().Local().TcpEndpoint, ":")
	rparts := strings.Split(node.TcpEndpoint, ":")
	node.TcpEndpoint = rparts[0] + ":" + rparts[1] + ":" + lparts[2]
	core.AppLog.Printf("Node : %v\n", node)
	return node
}
