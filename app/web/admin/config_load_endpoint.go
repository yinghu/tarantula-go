package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type ConfigLoader struct {
	*AdminService
}

func (s *ConfigLoader) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *ConfigLoader) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	if cid > 0 {
		conf, err := s.ItemService().LoadWithId(cid)
		if err != nil {
			session := core.OnSession{Successful: false, Message: err.Error()}
			w.Write(util.ToJson(session))
			return
		}
		w.Write(util.ToJson(conf))
		return
	}
	cname := r.PathValue("name")
	climt, err := strconv.ParseInt(r.PathValue("limit"), 10, 32)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	confs, err := s.ItemService().LoadWithName(cname, int(climt))
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.Write(util.ToJson(confs))
}
