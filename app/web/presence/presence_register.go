package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
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

		rf := event.RegisterEventFactory{}
		me := protocol.RegisterEvent{SystemId: uint64(login.SystemId), Name: login.Name, Source: "web"}
		tp, err := rf.FromRegisterEvent(&me)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to create topic %s", err.Error())
			return
		}
		tp.Event.Key.Array = core.ToBytes(s.Sequence())
		tp.NodeId = s.NodeId()
		tp.Tag = s.Context()
		_, err = s.Cluster().Publish(tp)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to publish topic %s", err.Error())
			return
		}
		mf := persistence.NewLoginObjectFactory()
		kv, err := mf.FromLoginObject(login)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to request %s", err.Error())
			return
		}

		tb := persistence.NewTaskBuilder(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "register"})

		jb1 := persistence.NewJobBuilder(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "validator"})

		jb1.Add(&protocol.Transaction{Meta: &protocol.Meta{Name: "register"}, Object: kv})

		jb := persistence.NewJobBuilder(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "job"})
		jb.Add(&protocol.Transaction{Meta: &protocol.Meta{Name: "grant"}, Object: kv})
		jb.Add(&protocol.Transaction{Meta: &protocol.Meta{Name: "update"}, Object: kv})
		tb.SetValidator(jb1.Target)
		tb.SetJob(jb.Target)
		rp, _ := s.Cluster().Issue(tb.Task())
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
