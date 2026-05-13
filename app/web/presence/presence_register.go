package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type PresenceRegister struct {
	*PresenceService
}

func (s *PresenceRegister) AccessControl() int32 {
	return core.PUBLIC_ACCESS_CONTROL
}

func (s *PresenceRegister) Register(login *protocol.LoginObject) (core.OnSession, error) {
	id, _ := s.Sequence().Id()
	login.SystemId = uint64(id)
	login.AccessControl = uint32(core.PROTECTED_ACCESS_CONTROL)
	hash, _ := s.Authenticator().HashPassword(login.Password)
	login.Password = hash
	err := s.SaveLogin(login)
	if err != nil {
		return core.OnSession{Successful: false, Message: err.Error()}, err
	}
	tk, err := s.Authenticator().CreateToken(int64(login.SystemId), int32(login.Id), int32(login.AccessControl))
	if err != nil {
		return core.OnSession{Successful: false, Message: err.Error()}, err
	}
	session := core.OnSession{Successful: true, SystemId: int64(login.SystemId), Stub: int32(login.Id), Token: tk, Home: s.F.Host}
	ticket, _ := s.Authenticator().CreateTicket(int64(login.SystemId), int32(login.Id), int32(login.AccessControl), bootstrap.TICKET_TIME_OUT_MINUTES)
	session.Ticket = ticket

	go func() {
		mf := persistence.NewCommodityObjectFactory()
		commodity := protocol.Commodity{Name: "gold", Type: "currency", TypeId: "hard currency", Amount: 12, Rechargeable: true}
		kv, err := mf.FromMessage(&commodity, mf.Header(persistence.COMMODITY_OBJECT_ID))
		kv.Key.Array = s.ToBytes(login.SystemId)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to request %s", err.Error())
			return
		}
		tb := persistence.NewTaskBuilder(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "register"})
		vb := tb.Validator(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "validator"})
		vb.Transaction().Meta(&protocol.Meta{Name: "register"}).Object(kv).Build()
		jb := tb.Job(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "job"})
		jb.Transaction().Meta(&protocol.Meta{Name: "grant"}).Object(kv).Build()
		jb.Build()
		rp, err := s.Cluster().Issue(tb.Build())
		if err != nil {
			core.AppLog.Debug().Msgf("TASK ERR %s", err.Error())
			return
		}
		core.AppLog.Debug().Msgf("TASK %v", rp)

	}()
	return session, nil
}

func (s *PresenceRegister) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	var login protocol.LoginObject
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	resp, _ := s.Register(&login)
	w.Write(util.ToJson(resp))
}
