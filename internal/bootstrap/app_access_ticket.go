package bootstrap

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AppAccessTicket struct {
	AppManager
	core.Ticket
}

func (s *AppAccessTicket) start() error {
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
	s.Ticket = &AuthManager{Tkn: &tkn, Cipher: &ci, Kid: "presence"}
	return nil
}
