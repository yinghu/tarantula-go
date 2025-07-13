package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminLoadCategory struct {
	*AdminService
}

func (s *AdminLoadCategory) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminLoadCategory) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pid, _ := strconv.Atoi(r.PathValue("id"))
	cid := int64(pid)
	if cid > 0 {
		cat, err := s.ItemService().LoadCategoryWithId(cid)
		if err != nil {
			session := core.OnSession{Successful: false, Message: err.Error()}
			w.Write(util.ToJson(session))
			return
		}
		w.Write(util.ToJson(cat))
		return
	}
	cname := r.PathValue("name")
	cat, err := s.ItemService().LoadCategory(cname)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.Write(util.ToJson(cat))
}
