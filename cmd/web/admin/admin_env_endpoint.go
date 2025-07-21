package main

import (
	"net/http"
	"os"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type WorkingEnv struct {
	CurBranch string `json:"CurBranch"`
}

type AdminEnv struct {
	*AdminService
}

func (s *AdminEnv) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminEnv) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	os.Chdir(s.publishDir)
	w.Write(util.ToJson(WorkingEnv{CurBranch: util.GitCurBranch().Message}))
}
