package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type PostofficeUnSubscriber struct {
	*PostofficeService
}

func (s *PostofficeUnSubscriber) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *PostofficeUnSubscriber) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tp event.SubscriptionEvent
	err := json.NewDecoder(r.Body).Decode(&tp)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	//id, err := s.createTopic(tp)
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: tp.App + "/" + tp.Name}))
}
