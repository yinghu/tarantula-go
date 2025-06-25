package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type SudoConfigNode struct {
	*AdminService
}

func (s *SudoConfigNode) AccessControl() int32 {
	return bootstrap.SUDO_ACCESS_CONTROL
}
func (s *SudoConfigNode) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var conf conf.Env
	json.NewDecoder(r.Body).Decode(&conf)
	s.Cluster().Atomic(conf.GroupName, func(ctx core.Ctx) error {
		return ctx.Put(conf.NodeName, string(util.ToJson(conf)))
	})
	w.Write(util.ToJson(conf))
}
