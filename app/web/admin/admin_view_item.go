package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminItemViewer struct {
	*AdminService
}

func (s *AdminItemViewer) AccessControl() int32 {
	return bootstrap.SUDO_ACCESS_CONTROL
}
func (s *AdminItemViewer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	conf, err := s.ItemService().InventoryManager().Load(cid)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(conf))
}
