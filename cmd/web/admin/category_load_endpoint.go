package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type CategoryLoader struct {
	*AdminService
}

func (s *CategoryLoader) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CategoryLoader) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	cname := r.PathValue("name")
	to := r.PathValue("to")
	target := r.PathValue("target")
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	if cid > 0 {
		cat, err := s.AdminService.ItemService().LoadCategoryWithId(cid)
		if err != nil {
			w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		}
		w.Write(util.ToJson(cat))
		return
	}
	if cname != "scope" {
		cat, err := s.AdminService.ItemService().LoadCategory(cname)
		if err != nil {
			w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		}
		w.Write(util.ToJson(cat))
		return
	}
	sp, err := strconv.ParseInt(to, 10, 32)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	list := s.ItemService().LoadCategories(int32(sp),target)
	w.Write(util.ToJson(list))
}
