package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
)

type AdminHashRingEndpoint struct {
	*AdminService
}

func (s *AdminHashRingEndpoint) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}

func (s *AdminHashRingEndpoint) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rq := make(chan []core.Node, 3)
	defer close(rq)
	s.Cluster().HashRing(core.RingRequest{Async: rq})
	n := <-rq
	data, err := json.Marshal(n)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
