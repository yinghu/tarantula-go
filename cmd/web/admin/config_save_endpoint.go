package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type ConfigSaver struct {
	*AdminService
}

func (s *ConfigSaver) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *ConfigSaver) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var conf item.Configuration
	err := json.NewDecoder(r.Body).Decode(&conf)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	sid, err := s.Sequence().Id()
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	conf.Id = sid
	err = s.ItemService().Validate(conf)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	err = s.ItemService().Save(conf)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	ch := make(chan core.OnSession, 1)
	defer close(ch)
	go s.PostJson("http://inventory:8080/inventory/itemadmin/saveconfig", conf, ch)
	ret := <-ch
	w.Write(util.ToJson(ret))
	//w.Write(util.ToJson(conf))
}
