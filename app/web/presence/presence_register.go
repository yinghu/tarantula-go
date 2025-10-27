package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type PresenceRegister struct {
	*PresenceService
}

func (s *PresenceRegister) AccessControl() int32 {
	return bootstrap.PUBLIC_ACCESS_CONTROL
}

func (s *PresenceRegister) Register(login bootstrap.Login) {
	id, _ := s.Sequence().Id()
	login.SystemId = id
	login.AccessControl = bootstrap.PROTECTED_ACCESS_CONTROL
	hash, _ := s.Authenticator().HashPassword(login.Hash)
	login.Hash = hash
	err := s.SaveLogin(&login)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.DB_OP_ERR_CODE)}
		return
	}
	tk, err := s.Authenticator().CreateToken(login.SystemId, login.Id, login.AccessControl)
	if err != nil {
		login.Cc <- core.Chunk{Remaining: false, Data: bootstrap.ErrorMessage(err.Error(), bootstrap.INVALID_TOKEN_CODE)}
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.Id, Token: tk, Home: s.Cluster().Local().HttpEndpoint}
	ticket, _ := s.AppAuth.CreateTicket(login.SystemId, login.Id, login.AccessControl,bootstrap.TICKET_TIME_OUT_SECONDS)
	session.Ticket = ticket
	login.Cc <- core.Chunk{Remaining: false, Data: util.ToJson(session)}
	go func() {
		id, err := s.Sequence().Id()
		if err != nil {
			return
		}
		me := event.RegisterEvent{SystemId: login.SystemId, Name: login.Name}
		me.OnOId(id)
		me.RegisterTime = time.Now()
		me.OnTopic("login")
		s.Send(&me)
		rw := item.OnInventory{SystemId: login.SystemId, ItemId: s.LoginReward.Id, Source: "login"}
		err = s.ItemService().InventoryManager().Grant(rw)
		if err != nil {
			core.AppLog.Printf("grant failed %s\n", err.Error())
		}
	}()
}

func (s *PresenceRegister) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	listener := make(chan core.Chunk)
	defer func() {
		close(listener)
		r.Body.Close()
	}()
	w.WriteHeader(http.StatusOK)
	var login bootstrap.Login
	json.NewDecoder(r.Body).Decode(&login)
	login.Cc = listener
	go s.Register(login)
	for c := range listener {
		w.Write(c.Data)
		if !c.Remaining {
			break
		}
	}

}
