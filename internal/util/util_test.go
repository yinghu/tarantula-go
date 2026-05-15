package util

import (
	"encoding/base64"
	"testing"
)

func TestPassword(t *testing.T) {
	h, err := HashPassword("password")
	if err != nil {
		t.Errorf("failed %s\n", err.Error())
	}
	er := ValidatePassword("password", h)
	if er != nil {
		t.Errorf("failed %s\n", er.Error())
	}
}

func TestPartition(t *testing.T) {
	p := Partition([]byte("hellp"), 5)
	if p > 5 {
		t.Errorf("falied partition %d\n", p)
	}
}

func TestKey(t *testing.T) {
	skey := KeyToBase64(Key(32))
	bkey, err := KeyFromBase64(skey)
	if err != nil {
		t.Errorf("bad format key %s\n", err.Error())
	}
	ckey := base64.StdEncoding.EncodeToString(bkey)
	if skey != ckey {
		t.Errorf("key not same %s %s \n", skey, ckey)
	}
}


func TestVaultClient(t *testing.T) {
	vclient := VaultClient{Host: "http://192.168.1.11:8200", Token: ""}
	err := vclient.Auth()
	if err != nil {
		t.Errorf("error %s", err.Error())
	}
	sk, err := vclient.GetSecret("dev/presence", "gcp")
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	iam, _ := sk.Data["iam"].(string)
	gcp := GcpComputeEngine{
		ServiceAccount: iam,
		ProjectId:      "prismatic-grail-206205",
		Zone:           "us-east1-c",
	}
	err = gcp.Auth()
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	defer gcp.Close()
	//err = gcp.Insert("tarantula-build-01")
	//if err != nil {
	//t.Errorf("error %s", err.Error())
	//}
	ins, err := gcp.Get("tarantula-build-01")
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	pk, _ := sk.Data["ssh"].(string)
	ssh := SshClient{Host: ins.GetNetworkInterfaces()[0].AccessConfigs[0].GetNatIP(), User: "yinghu_lu", PrivateKey: pk}
	err = ssh.WithKey()
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	defer ssh.Close()
	ssh.Run("pwd && git --version && docker --version && df -h")
	//fmt.Printf("Key %v", sk.Data)
}
