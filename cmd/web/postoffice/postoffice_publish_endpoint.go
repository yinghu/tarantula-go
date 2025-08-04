package main

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	topic := r.PathValue("topic")
	cid, err := strconv.ParseInt(r.PathValue("cid"), 10, 64)
	if err!=nil{
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	me, err := s.Create(int(cid),"ticket")
	if err!=nil{
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	err = json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: topic}))
	view := s.Cluster().View()
	for i := range view {
		v := view[i]
		core.AppLog.Printf("Sending to : %s,%s,%s %d\n", v.Name, v.TcpEndpoint, s.Cluster().Local().Name,cid) // no prefix
		if v.Name == s.Cluster().Local().Name {
			continue
		}
		pub := event.SocketPublisher{Remote: v.TcpEndpoint}
		pub.Publish(me, "ticket")
	}
}
