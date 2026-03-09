package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminCreateAccessKey struct {
	*AdminService
}

type KeyExpiration struct {
	ExpiryTime time.Time `json:"ExpiryTime"`
	Key        string    `json:"AccessKey"`
}

func (s *AdminCreateAccessKey) AccessControl() int32 {
	return core.ADMIN_ACCESS_CONTROL
}
func (s *AdminCreateAccessKey) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	rq := make(chan core.Chunk, 10)
	defer close(rq)
	creq := core.DataRequest{Opt: core.CREATE_DATA_REQUEST, Key: []byte("key90"), Value: []byte("value2"), Async: rq}
	s.Cluster().Request(creq)
	for chunk := range rq {
		core.AppLog.Debug().Msgf("Chunk %s", string(chunk.Data))
		if !chunk.Remaining {
			break
		}
	}
	defer r.Body.Close()
	var cp KeyExpiration
	err := json.NewDecoder(r.Body).Decode(&cp)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	if cp.ExpiryTime.Before(time.Now()) {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: fmt.Sprintf("time already expired :%v", cp.ExpiryTime)}))
		return
	}
	dur := int(time.Until(cp.ExpiryTime).Seconds())
	key, err := s.Authenticator().CreateTicket(rs.SystemId, rs.Stub, rs.AccessControl, dur)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	cp.Key = key
	w.Write(util.ToJson(cp))

}
