package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminPublisher struct {
	*AdminService
}

func (s *AdminPublisher) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminPublisher) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	env := r.PathValue("repo")
	gr := util.GitPush()
	core.AppLog.Printf("Publish repo : %s\n", env)
	w.Write(util.ToJson(gr))
}
