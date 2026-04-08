package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type PresenceLogin struct {
	*PresenceService
}

func (s *PresenceLogin) AccessControl() int32 {
	return core.PUBLIC_ACCESS_CONTROL
}

func (s *PresenceLogin) Login(login bootstrap.Login) {
	pwd := login.Hash
	err := s.LoadLogin(&login)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	err = util.ValidatePassword(pwd, login.Hash)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.WRONG_PASS_CODE)}
		return
	}
	tk, err := s.Authenticator().CreateToken(login.SystemId, login.Id, login.AccessControl)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.INVALID_TOKEN_CODE)}
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.Id, Token: tk, Home: s.F.Host}
	ticket, _ := s.Authenticator().CreateTicket(login.SystemId, login.Id, login.AccessControl, bootstrap.TICKET_TIME_OUT_MINUTES)
	session.Ticket = ticket
	login.Cc <- core.Chunk{Remaining: false, Data: util.ToJson(session)}
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
		rw := item.OnInventory{SystemId: login.SystemId, ItemId: s.LoginReward.Id, Source: "login"}
		err = s.ItemService().InventoryManager().Grant(rw)
		if err != nil {
			core.AppLog.Printf("grant failed %s\n", err.Error())
		}
	}()
}

func (s *PresenceLogin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	listener := make(chan core.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	var login bootstrap.Login
	json.NewDecoder(r.Body).Decode(&login)
	login.Cc = listener
	go s.Login(login)
	for c := range listener {
		cv, _ := c.Data.([]byte)
		w.Write(cv)
		if !c.Remaining {
			break
		}
	}

}
