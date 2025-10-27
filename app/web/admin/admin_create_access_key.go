package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminCreateAccessKey struct {
	*AdminService
}

type KeyExpiration struct {
	ExpiryTime time.Time `json:"ExpiryTime"`
}

func (s *AdminCreateAccessKey) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminCreateAccessKey) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var cp KeyExpiration
	json.NewDecoder(r.Body).Decode(&cp)
	dur := int(time.Until(cp.ExpiryTime).Seconds())
	key, err := s.AppAuth.CreateTicket(rs.SystemId, rs.Stub, rs.AccessControl, dur)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	session := core.OnSession{Successful: true, Message: fmt.Sprintf("%s : %d", key, dur)}
	w.Write(util.ToJson(session))
}
