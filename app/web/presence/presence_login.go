package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type PresenceLogin struct {
	*PresenceService
}

func (s *PresenceLogin) AccessControl() int32 {
	return core.PUBLIC_ACCESS_CONTROL
}

func (s *PresenceLogin) Login(login *protocol.LoginObject) (core.OnSession, error) {
	pwd := login.Password
	err := s.LoadLogin(login)
	if err != nil {
		return core.OnSession{Successful: false, Message: err.Error()}, err
	}
	err = util.ValidatePassword(pwd, login.Password)
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
		id, err := s.Sequence().Id()
		if err != nil {
			core.AppLog.Warn().Msgf("failed to create seq %s", err.Error())
			return
		}
		rf := event.LoginEventFactory{}
		me := protocol.LoginEvent{SystemId: uint64(login.SystemId), Name: login.Name, Source: "web"}
		tp, err := rf.FromLoginEvent(&me)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to create topic %s", err.Error())
			return
		}
		tp.Event.Id = uint64(id)
		tp.NodeId = s.NodeId()
		tp.Tag = s.Context()
		_, err = s.Cluster().Publish(tp)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to publish topic %s", err.Error())
			return
		}
		//rw := item.OnInventory{SystemId: login.SystemId, ItemId: s.LoginReward.Id, Source: "login"}
		//err = s.ItemService().InventoryManager().Grant(rw)
		//if err != nil {
		//core.AppLog.Printf("grant failed %s\n", err.Error())
		//}
	}()
	return session, nil
}

func (s *PresenceLogin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	var login protocol.LoginObject
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	resp, _ := s.Login(&login)
	w.Write(util.ToJson(resp))

}
