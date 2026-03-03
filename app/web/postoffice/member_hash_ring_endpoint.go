package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
)

type MemberHashRingEndpoint struct {
	*PostofficeService
}

func (s *MemberHashRingEndpoint) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *MemberHashRingEndpoint) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rq := make(chan []core.Node, 1)
	defer close(rq)
	s.mm.HashRing(core.RingRequest{Async: rq})
	n := <-rq
	data, err := json.Marshal(n)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
