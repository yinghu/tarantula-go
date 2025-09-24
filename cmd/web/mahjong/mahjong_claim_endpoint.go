package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type ClaimHand struct {
	ClaimList string `json:"HandList"`
}

type MahjongClaimer struct {
	*MahjongService
}

func (s *MahjongClaimer) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *MahjongClaimer) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ch := ClaimHand{}
	err := json.NewDecoder(r.Body).Decode(&ch)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	h := Hand{}
	h.New()
	for c := range strings.SplitSeq(ch.ClaimList, ",") {
		t := Tile{}
		t.From(c)
		h.Tiles = append(h.Tiles, t)
	}
	claimed := s.Mahjong(&h)
	if !claimed {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: "not claimed"}))
		return
	}
	w.Write(util.ToJson(h.Formed))

}
