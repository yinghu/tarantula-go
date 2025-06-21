package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminLoadConfig struct {
	*AdminService
}

func (s *AdminLoadConfig) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminLoadConfig) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cname := r.PathValue("name")
	climt, err := strconv.Atoi(r.PathValue("limit"))
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	confs, err := s.ItemService().LoadWithName(cname, climt)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.Write(util.ToJson(confs))
}
