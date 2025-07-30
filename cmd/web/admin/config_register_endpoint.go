package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type ConfigRegister struct {
	*AdminService
}

func (s *ConfigRegister) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *ConfigRegister) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	opt := r.PathValue("opt")
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	app := r.PathValue("app")
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	switch opt {
	case "save":
		var reg item.ConfigRegistration
		err := json.NewDecoder(r.Body).Decode(&reg)
		if err != nil {
			w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
			return
		}
		app := reg.App
		if !slices.Contains(s.managedApps, app) {
			w.Write(util.ToJson(core.OnSession{Successful: false, Message: "app not existed"}))
			return
		}
		err = s.AdminService.ItemService().Register(reg)
		if err != nil {
			w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
			return
		}
		w.Write(util.ToJson(core.OnSession{Successful: true, Message: "registered"}))
	case "load":
		reg, err := s.ItemService().Check(cid, app)
		if err == nil {
			w.Write(util.ToJson(reg))
		} else {
			w.Write(util.ToJson(reg))
		}
	case "delete":
		err = s.ItemService().Release(int32(cid))
		if err != nil {
			w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		} else {
			w.Write(util.ToJson(core.OnSession{Successful: true, Message: "deleted"}))
		}
	default:
		w.Write(util.ToJson(core.OnSession{Successful: true, Message: "not supported"}))
	}
}
