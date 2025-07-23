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
	session := core.OnSession{Successful: true, Message: "app cluster admin [" + s.Cluster().Group() + "]"}
	if cmd == "join" {
		var join core.Node
		json.NewDecoder(r.Body).Decode(&join)
		s.Cluster().OnJoin(s.convert(join))
		w.WriteHeader(http.StatusOK)
		w.Write(util.ToJson(session))
		return
	}
	if cmd == "left" {
		var left core.Node
		json.NewDecoder(r.Body).Decode(&left)
		s.Cluster().OnLeave(s.convert(left))
		w.WriteHeader(http.StatusOK)
		w.Write(util.ToJson(session))
		return
	}
	if cmd == "update" {
		var update KVUpdate
		json.NewDecoder(r.Body).Decode(&update)
		s.Cluster().OnUpdate(update.Key, update.Value, update.Opt)
		w.WriteHeader(http.StatusOK)
		w.Write(util.ToJson(session))
		s.ItemListener().OnUpdate()
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
	return node
}
