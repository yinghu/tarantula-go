package main

import (
	"net/http"
	"os"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type CategoryPublisher struct {
	*AdminService
}

func (s *CategoryPublisher) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CategoryPublisher) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	env := r.PathValue("env")
	sid := r.PathValue("id")
	cid, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	conf, err := s.ItemService().LoadCategoryWithId(cid)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	dest, err := os.OpenFile(s.publishDir+"/"+sid+".json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	defer dest.Close()
	dest.WriteString(string(util.ToJson(conf)))
	os.Chdir(s.publishDir)
	gr := util.GitPull()
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitAdd(sid + ".json")
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitCommit("publish config :" + sid + ".json")
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	gr = util.GitPush()
	if !gr.Successful {
		w.Write(util.ToJson(gr))
		return
	}
	core.AppLog.Printf("Publish category :%d %s\n", cid, env)
	w.Write(util.ToJson(util.GitStatus()))
}
