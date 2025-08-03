package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type PostofficePublisher struct {
	*PostofficeService
}

func (s *PostofficePublisher) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *PostofficePublisher) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	key := r.PathValue("key")
	var me event.MessageEvent
	err := json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: key}))
	view := s.Cluster().View()
	for i := range view {
		v := view[i]
		if v.Name == s.Cluster().Local().Name {
			continue
		}
		pub := event.SocketPublisher{Remote: v.TcpEndpoint}
		pub.Publish(&me, "ticket")
	}
}
