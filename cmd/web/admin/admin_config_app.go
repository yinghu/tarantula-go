package main

import (
	"fmt"
	"io"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminConfigApp struct {
	*AdminService
	core.Ticket
}

func (s *AdminConfigApp) start() error {
	tkn := util.JwtHMac{Alg: core.JWT_ALG, Ksz: core.JWT_KEY_SIZE}
	ci := util.Aes{Ksz: core.CIPHER_KEY_SIZE}
	err := s.Cluster().Atomic("presence", func(ctx core.Ctx) error {
		jsk, err := ctx.Get(core.JWT_KEY_NAME)
		if err != nil {
			fmt.Println("Create new jwt key")
			nkey := util.Key(tkn.Ksz)
			ctx.Put(core.JWT_KEY_NAME, util.KeyToBase64(nkey))
			tkn.HMacFromKey(nkey)
			return nil
		}
		jk, err := util.KeyFromBase64(jsk)
		if err != nil {
			return err
		}
		tkn.HMacFromKey(jk)
		return nil
	})
	if err != nil {
		return err
	}
	err = s.Cluster().Atomic("presence", func(ctx core.Ctx) error {
		csk, err := ctx.Get(core.CIPHER_KEY_NAME)
		if err != nil {
			fmt.Println("Create new cipher key")
			ckey := util.Key(ci.Ksz)
			ctx.Put(core.CIPHER_KEY_NAME, util.KeyToBase64(ckey))
			ci.AesGcmFromKey(ckey)
		}
		ck, err := util.KeyFromBase64(csk)
		if err != nil {
			return err
		}
		ci.AesGcmFromKey(ck)
		return nil
	})
	if err != nil {
		return err
	}
	s.Ticket = &bootstrap.AuthManager{Tkn: &tkn, Cipher: &ci, Kid: "presence"}
	return nil
}

func (s *AdminConfigApp) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *AdminConfigApp) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	app := r.PathValue("app")
	tick, err := s.CreateTicket(rs.SystemId, rs.Stub, rs.AccessControl)
	if err != nil {
		session := core.OnSession{Successful: true, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://web/"+app+"/admin", nil)
	if err != nil {
		session := core.OnSession{Successful: true, Message: err.Error()}
		w.Write(util.ToJson(session))
	}
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		session := core.OnSession{Successful: true, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}
