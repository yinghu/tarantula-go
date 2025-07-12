package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminLoadCategories struct {
	*AdminService
}

func (s *AdminLoadCategories) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminLoadCategories) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cname := r.PathValue("scope")
	list := s.ItemService().FromScope(cname)
	w.Write(util.ToJson(list))
}
