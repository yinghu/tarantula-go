package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
)

type AdminKeyRingEndpoint struct {
	*AdminService
}

func (s *AdminKeyRingEndpoint) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}

func (s *AdminKeyRingEndpoint) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	key := r.PathValue("key")
	rq := make(chan []core.Node, 1)
	defer close(rq)
	s.Cluster().HashRing(core.RingRequest{Async: rq, Opt: core.REPLICA_RING_OPT, Token: s.Cluster().RingToken([]byte(key))})
	n := <-rq
	data, err := json.Marshal(n)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
