package main

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminParseSnowFlakeId struct {
	*AdminService
}

type SnowFlakeResp struct {
	Id        int64 `json:"id"`
	NodeId    int64 `json:"nodeId"`
	Timestamp int64 `json:"timestamp"`
	Sequence  int64 `json:"seqence"`
}

func (s *AdminParseSnowFlakeId) AccessControl() int32 {
	return bootstrap.SUDO_ACCESS_CONTROL
}
func (s *AdminParseSnowFlakeId) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	if cid <= 0 {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: "id cannot be less 0"}))
		return
	}
	t,n,se := s.Sequence().Parse(cid)
	ps := SnowFlakeResp{Timestamp: t,NodeId: n,Sequence:se,Id: cid}
	w.Write(util.ToJson(ps))
}
