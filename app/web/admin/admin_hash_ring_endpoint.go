package main

import (
	"io"
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminHashRingEndpoint struct {
	*AdminService
}

func (s *AdminHashRingEndpoint) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}

func (s *AdminHashRingEndpoint) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	stream, err := s.Cluster().HashRing(core.RingRequest{})
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	ring := make([]core.Node, 0)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("streaming error %s", err.Error())
			break
		}
		ring = append(ring, core.Node{Name: data.Name, RingToken: data.Hash, RpcEndpoint: data.Endpoint, IP: data.Address})
	}
	w.Write(util.ToJson(ring))
}
