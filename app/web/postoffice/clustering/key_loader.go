package clustering

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

type KeyLoader struct {
	config core.Vault
	vault  util.VaultClient
}

func (a *KeyLoader) auth() error {
	a.vault = util.VaultClient{Host: a.config.Host, Token: a.config.Token}
	return a.vault.Auth()
}

func (a *KeyLoader) load(mount string, path string) (*protocol.AuthKey, error) {
	ak := protocol.AuthKey{}
	kv, err := a.vault.GetSecret(mount, path)
	if err != nil {
		core.AppLog.Warn().Msgf("load auth key error %s", err.Error())
		return &ak, err
	}
	jwt, _ := kv.Data["jwt"].(string)
	cipher, _ := kv.Data["cipher"].(string)

	ak.Jwt = []byte(jwt)       //[]byte("M4KExRr9HQbFRmUObhPsFJ2aDC4e/dkayHqwA568czw=")
	ak.Cipher = []byte(cipher) //("NyXPfuLkqlsZcz//Y7Cm7oQ14MgIv29ouRmia+Vr2NY=")

	return &ak, nil
}
