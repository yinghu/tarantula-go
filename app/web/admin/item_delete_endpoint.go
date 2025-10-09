package main

import (
	"errors"
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type ItemDeleter struct {
	*AdminService
}

func (s *ItemDeleter) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *ItemDeleter) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tp := r.PathValue("type")
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	switch tp {
	case "category":
		err = s.AdminService.ItemService().DeleteCategoryWithId(cid)
	case "config":
		err = s.AdminService.ItemService().DeleteWithId(cid)
	default:
		err = errors.New("type not supported :" + tp)
	}
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "deleted"}))
}
