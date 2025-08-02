package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
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
	key := r.PathValue("key")
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: key}))
}
