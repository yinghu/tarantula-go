package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type CSQueryer struct {
	*AdminService
}

func (s *CSQueryer) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}

func (s *CSQueryer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	qid, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	me := event.CreateQuery(uint32(qid))
	err = json.NewDecoder(r.Body).Decode(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	defer close(me.QCc())
	go s.List(me)
	//w.Write([]byte("["))
	for c := range me.QCc() {
		if !c.Remaining {
			break
		}
		resp, ok := c.Data.(*protocol.Response)
		if ok {
			for _, data := range resp.Data.List {
				me := event.MessageEvent{}
				core.Import(&me, data.Key, data.Value, 200)
				core.AppLog.Debug().Msgf("Data : %v", me)
				//w.Write(util.ToJson(me))

			}
		}
	}
	//w.Write([]byte("]"))
	w.Write(util.ToJson(me))
}
