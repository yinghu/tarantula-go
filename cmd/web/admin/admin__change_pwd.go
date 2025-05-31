package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type AdminChangePwd struct {
	*AdminService
	accessControl int32
}

func (s AdminChangePwd) Login(login *event.Login) error {

	return nil
}
func (s *AdminChangePwd) AccessControl() int32 {
	return s.accessControl
}
func (s *AdminChangePwd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tkn := r.Header.Get("Authorization")
	fmt.Printf("TX : %s\n", tkn)
	parts := strings.Split(tkn, " ")
	sess, err := s.Auth.ValidateToken(parts[1])
	if err != nil {
		fmt.Printf("Err : %s\n", err.Error())
	}
	fmt.Printf("Sess: %d\n", sess.SystemId)
	var login event.Login
	json.NewDecoder(r.Body).Decode(&login)
	pwd := login.Hash
	err = s.LoadLogin(&login)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	err = s.Auth.ValidatePassword(pwd, login.Hash)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	tk, err := s.Auth.CreateToken(login.SystemId, login.SystemId, 1)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	session := core.OnSession{Successful: true, SystemId: login.SystemId, Stub: login.SystemId, Token: tk, Home: s.Cluster().Local().HttpEndpoint}
	w.Write(util.ToJson(session))
}
