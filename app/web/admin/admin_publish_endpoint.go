package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type AdminPublisher struct {
	*AdminService
}

func (s *AdminPublisher) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminPublisher) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	repo := item.RepoUpdate{Admin: s.Cluster().Group()}
	defer r.Body.Close()
	defer func() {
		s.Cluster().Atomic(s.Cluster().Group(), func(ctx core.Ctx) error {
			ctx.Put("push", string(util.ToJson(repo)))
			return nil
		})
	}()
	err := json.NewDecoder(r.Body).Decode(&repo)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.WriteHeader(http.StatusOK)
	cur := util.GitCurBranch().Message
	repo.Source = cur
	gr := util.GitPush()
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	core.AppLog.Printf("Publish repo : %s : %s\n", repo.Target, cur)
	if cur == repo.Target {
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitCheckout(repo.Target)
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
