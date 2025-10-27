package main

import (
	"encoding/json"
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
	session := core.OnSession{Successful: true, Message: "key created"}
	w.Write(util.ToJson(session))
}
