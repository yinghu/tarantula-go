package main

import (
	"encoding/json"
	"net/http"

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
	var reg item.ConfigRegistration
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	err = s.AdminService.ItemService().Register(reg)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "registered"}))
}
