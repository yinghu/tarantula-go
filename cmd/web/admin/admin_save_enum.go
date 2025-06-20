package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type AdminSaveEnum struct {
	*AdminService
}

func (s *AdminSaveEnum) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminSaveEnum) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var conf item.Enum
	err := json.NewDecoder(r.Body).Decode(&conf)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	err = s.ItemService().SaveEnum(conf)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.Write(util.ToJson(conf))
}
