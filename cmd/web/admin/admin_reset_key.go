package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminResetKey struct {
	*AdminService
}

func (s *AdminResetKey) AccessControl() int32 {
	return bootstrap.SUDO_ACCESS_CONTROL
}
func (s *AdminResetKey) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	s.Cluster().Atomic("presence", func(ctx cluster.Ctx) error {
		//ctx.Get(core.JWT_KEY_NAME)
		ctx.Put("teset", "test key")
		return nil
	})
	session := core.OnSession{Successful: true, Message: "admin reset key"}

	w.Write(util.ToJson(session))
}
