package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type PresenceAdmin struct {
	*PresenceService
}

func (s *PresenceAdmin) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *PresenceAdmin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	session := core.OnSession{Successful: true, Message: "presence admin"}
	w.Write(util.ToJson(session))
}
