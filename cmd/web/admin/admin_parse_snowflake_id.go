package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminParseSnowFlakeId struct {
	*AdminService
}

type SnowFlakeResp struct {
	Id        int64 `json:"snowFlakeId,string"`
	NodeId    int64 `json:"nodeId,string"`
	Timestamp int64 `json:"timestamp,string"`
	Sequence  int64 `json:"seqence,string"`
}

func (s *AdminParseSnowFlakeId) AccessControl() int32 {
	return bootstrap.SUDO_ACCESS_CONTROL
}
func (s *AdminParseSnowFlakeId) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var e SnowFlakeResp
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	w.WriteHeader(http.StatusOK)

	if e.Id <= 0 {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: "id cannot be less 0"}))
		return
	}
	t, n, se := s.Sequence().Parse(e.Id)
	ps := SnowFlakeResp{Timestamp: t, NodeId: n, Sequence: se, Id: e.Id}
	w.Write(util.ToJson(ps))
}
