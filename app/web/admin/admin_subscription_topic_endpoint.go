package main

import (
	"io"
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type AdminSubscriptionTopicEndpoint struct {
	*AdminService
}

func (s *AdminSubscriptionTopicEndpoint) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}
func (s *AdminSubscriptionTopicEndpoint) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	subs := make([]*protocol.Subscription, 0)
	stream, err := s.Cluster().TopicList()
	if err != nil {
		w.Write(util.ToJson(subs))
		return
	}
	for {
		sub, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			core.AppLog.Debug().Msgf("streaming error %s", err.Error())
			break
		}
		subs = append(subs, sub)
	}
	w.Write(util.ToJson(subs))
}
