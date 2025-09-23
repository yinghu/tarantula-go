package main

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	//"gameclustering.com/internal/util"
)

type ClaimHand struct {
	ClaimList string `json:"HandList,string"`
}

type MahjongClaimer struct {
	*MahjongService
}

func (s *MahjongClaimer) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *MahjongClaimer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	h := ClaimHand{}
	err := json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	w.Write(util.ToJson(h))
}
