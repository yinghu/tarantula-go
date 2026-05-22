package util

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	Host   string
	Token  string
	client *vault.Client
}

func (v *VaultClient) Auth() error {
	cfg := vault.DefaultConfig()
	cfg.Address = v.Host
	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}
	client.SetToken(v.Token)
	v.client = client
	return nil
}

func (v *VaultClient) GetSecret(mountPath string, path string) (*vault.KVSecret, error) {
	return v.client.KVv2(mountPath).Get(context.Background(), path)
}

func (v *VaultClient) SetSecret(mountPath string, path string, value string) (*vault.KVSecret, error) {
	kv := make(map[string]any)
	kv[path] = value
	return v.client.KVv2(mountPath).Put(context.Background(), mountPath, kv)
}

func (v *VaultClient) CreateStore(){
	mount := vault.MountInput{Type: "kv",Options:map[string]string{},Description: ""}
	v.client.Sys().Mount("",&mount)
	
}
