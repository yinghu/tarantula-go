package util

import (
	"encoding/base64"
	"fmt"
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

func TestSsh(t *testing.T) {
	sh := SshClient{
		Host:     "192.168.1.11",
		User:     "yinghu",
		Password: "casino123",
	}
	err := sh.WithPassword()
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	defer sh.Close()
	sh.Run("cd ~/development/icodesoftware/tarantula-go && git pull")
}

func TestGcpAuth(t *testing.T) {
	gcp := GcpComputeEngine{
		ServiceAccount: "C:\\development\\prismatic-grail-206205-bdc198bbb5de.json",
		ProjectId:      "prismatic-grail-206205",
		Zone:           "us-east1-c",
	}
	err := gcp.Auth()
	if err != nil {
		t.Errorf("error %s", err.Error())
	}
	defer gcp.Close()

	//err = gcp.Insert("tarantula-build-01")
	//if err != nil {
	//t.Errorf("error %s", err.Error())
	//}
	//ins, err := gcp.Get("tarantula-build-01")
	//if err != nil {
	//t.Errorf("error %s", err.Error())
	//return
	//}
	//ssh := SshClient{Host: ins.GetNetworkInterfaces()[0].AccessConfigs[0].GetNatIP(), User: "yinghu_lu", PrivateKey: "C:\\Users\\yingh\\.ssh\\gx-01.key"}
	//err = ssh.WithKey()
	//if err != nil {
	//t.Errorf("error %s", err.Error())
	//}
	//defer ssh.Close()
	//ssh.Run("pwd && git --version && docker --version && df -h")
	//f, err := os.Open("C:\\development\\s3.json")
	//if err != nil {
	//t.Errorf("error %s", err.Error())
	//}
	//err = ssh.Upload(*f, "/home/yinghu_lu/s3.json", "0544")

	//err = gcp.Delete("tarantula-build-01")
	//if err != nil {
	//t.Errorf("error %s", err.Error())
	//}
}

func TestVaultClient(t *testing.T) {
	vclient := VaultClient{Host: "http://192.168.1.11:8200", Token: ""}
	err := vclient.Auth()
	if err != nil {
		t.Errorf("error %s", err.Error())
	}
	sk, err := vclient.GetSecret("tarantula", "gcp")
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	fmt.Printf("Key %v", sk.Data)
}
