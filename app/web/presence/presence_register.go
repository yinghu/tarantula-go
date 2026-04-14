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
		id, err := s.Sequence().Id()
		if err != nil {
			core.AppLog.Warn().Msgf("failed to create seq %s", err.Error())
			return
		}
		rf := event.RegisterEventFactory{}
		me := protocol.RegisterEvent{SystemId: uint64(login.SystemId), Name: login.Name, Source: "web"}
		tp, err := rf.FromRegisterEvent(&me)
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
		mf := persistence.NewLoginObjectFactory()
		kv, err := mf.FromLoginObject(login)
		if err != nil {
			core.AppLog.Warn().Msgf("failed to request %s", err.Error())
			return
		}
		//req, err := mf.Request(kv)
		//if err != nil {
		//core.AppLog.Warn().Msgf("failed to request %s", err.Error())

		//return
		//}
		//resp, err := s.Cluster().Request(req)
		//if err != nil {
		//core.AppLog.Warn().Msgf("failed to request %s", err.Error())
		//return
		//}
		//core.AppLog.Debug().Msgf("REQ %v", resp)
		tsk := protocol.Task{NodeId: s.NodeId(), Tag: s.Context(), Name: "register", Prefix: s.Cluster().RingToken([]byte(login.Name))}

		//obj, err := anypb.New(login)
		//if err != nil {
		//return
		//}
		ts := make([]*protocol.Transaction, 0)
		ts = append(ts, &protocol.Transaction{Meta: &protocol.Meta{Name: "register"}, Object: kv})
		tsk.Transactions = ts
		rp, _ := s.Cluster().Issue(&tsk)
		core.AppLog.Debug().Msgf("TASK %v", rp)

		//req := mf.Request()
		//rw := item.OnInventory{SystemId: login.SystemId, ItemId: s.LoginReward.Id, Source: "login"}
		//err = s.ItemService().InventoryManager().Grant(rw)
		//if err != nil {
		//core.AppLog.Printf("grant failed %s\n", err.Error())
		//}
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
