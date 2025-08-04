package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type PostofficeSubscriber struct {
	*PostofficeService
}

func (s *PostofficeSubscriber) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *PostofficeSubscriber) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	app := r.PathValue("app")
	topic := r.PathValue("topic")
	pt := Topic{Name: topic, App: app}
	id, err := s.createTopic(pt)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	pt.Id = id
	s.topics[id] = pt
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: app + "/" + topic}))
}
