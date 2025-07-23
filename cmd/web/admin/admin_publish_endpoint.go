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
	cur := util.GitCurBranch().Message

	gr := util.GitPush()
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	core.AppLog.Printf("Publish repo : %s : %s\n", env, cur)
	if cur == env {
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitCheckout(env)
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitPull()
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitMerge(cur)
	if !gr.Successful {
		util.GitCheckout(cur)
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitPush()
	util.GitCheckout(cur)
	w.Write(util.ToJson(gr))
}
