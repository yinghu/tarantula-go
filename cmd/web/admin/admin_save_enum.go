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
	err = s.ItemService().ValidateEnum(conf)
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
	ch := make(chan core.OnSession, 1)
	defer close(ch)
	go s.PostJson("http://inventory:8080/inventory/itemadmin/saveenum", conf, ch)
	ret := <-ch
	w.Write(util.ToJson(ret))
}
