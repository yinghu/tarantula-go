package main

import (
	"io"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminConfigApp struct {
	*AdminService
}

func (s *AdminConfigApp) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminConfigApp) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	app := r.PathValue("app")
	tick, err := s.AppAuth.CreateTicket(rs.SystemId, rs.Stub, rs.AccessControl)
	if err != nil {
		session := core.OnSession{Successful: true, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://web/"+app+"/admin/test", nil)
	if err != nil {
		session := core.OnSession{Successful: true, Message: err.Error()}
		w.Write(util.ToJson(session))
	}
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		session := core.OnSession{Successful: true, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
