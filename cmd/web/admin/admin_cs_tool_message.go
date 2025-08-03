package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type CSMessager struct {
	*AdminService
}

func (s *CSMessager) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CSMessager) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "message"}))
}
