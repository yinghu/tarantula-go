package main

import (
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"
)

type PresenceClusterGet struct {
	*PresenceService
}

func (s *PresenceClusterGet) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *PresenceClusterGet) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	kn := r.PathValue("key")
	login := protocol.LoginObject{Name: kn}
	mf := persistence.NewLoginObjectFactory()
	kv, err := mf.FromLoginObject(&login)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	req, err := mf.Request(kv)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	req.Opt = core.GET_DATA_REQUEST
	resp, err := s.Cluster().Request(req)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}

	if err := mf.Set(resp).QList(func(h *protocol.Header, m proto.Message) bool {
		w.Write(util.ToJson(m))
		return false
	}); err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: "not found"}))
	}
}
