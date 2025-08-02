package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
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
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: key}))
}
