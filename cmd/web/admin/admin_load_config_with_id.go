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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	conf, err := s.ItemService().LoadWithId(int64(id))
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.Write(util.ToJson(conf))
}
