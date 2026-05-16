package clustering

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	vault "github.com/hashicorp/vault/api"
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
	switch path {
	case "auth":
		return a.toAuthKey(kv), nil
	case "sql":
		return a.toSqlKey(kv), nil
	case "gcp":
		return a.toGcpKey(kv), nil
	case "git":
		return a.toGitKey(kv), nil
	}
	return &ak, fmt.Errorf("key path not existed")
}

func (a *KeyLoader) toAuthKey(kv *vault.KVSecret) *protocol.AuthKey {
	jwt, _ := kv.Data["jwt"].(string)
	cipher, _ := kv.Data["cipher"].(string)
	ak := protocol.AuthKey{}
	ak.Jwt = []byte(jwt)    
	ak.Cipher = []byte(cipher)
	return &ak
}

func (a *KeyLoader) toSqlKey(kv *vault.KVSecret) *protocol.AuthKey {
	user, _ := kv.Data["user"].(string)
	url, _ := kv.Data["url"].(string)
	ak := protocol.AuthKey{Sql: &protocol.SqlAccess{}}
	ak.Sql.User = user
	ak.Sql.Url = url
	return &ak
}

func (a *KeyLoader) toGcpKey(kv *vault.KVSecret) *protocol.AuthKey {
	iam, _ := kv.Data["iam"].(string)
	ssh, _ := kv.Data["ssh"].(string)
	ak := protocol.AuthKey{Gcp: &protocol.GcpAccess{}}
	ak.Gcp.Iam = iam
	ak.Gcp.Ssh = ssh
	return &ak
}

func (a *KeyLoader) toGitKey(kv *vault.KVSecret) *protocol.AuthKey {
	key, _ := kv.Data["key"].(string)
	user, _ := kv.Data["user"].(string)
	email, _ := kv.Data["email"].(string)
	ak := protocol.AuthKey{Git: &protocol.GitAccess{}}
	ak.Git.Key = key
	ak.Git.User = user
	ak.Git.Email = email
	return &ak
}
