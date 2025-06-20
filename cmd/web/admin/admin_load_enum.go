package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminLoadEnum struct {
	*AdminService
}

func (s *AdminLoadEnum) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminLoadEnum) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cname := r.PathValue("name")
	cat, err := s.ItemService().LoadEnum(cname)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.Write(util.ToJson(cat))
}
