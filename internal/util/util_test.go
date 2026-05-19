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

func TestVaultClient(t *testing.T) {
	vclient := VaultClient{Host: "http://192.168.1.11:8200", Token: ""}
	err := vclient.Auth()
	if err != nil {
		t.Errorf("error %s", err.Error())
	}
	sk, err := vclient.GetSecret("dev/presence", "git")
	if err != nil {
		t.Errorf("error %s", err.Error())
		return
	}
	tok, _ := sk.Data["token"].(string)
	org, _ := sk.Data["org"].(string)
	gh := GitHubApi{Token: tok, Org: org}
	repos, err := gh.ListRepos()
	if err != nil {
		t.Errorf("err %s", err.Error())
		return
	}
	fmt.Printf("REPO %v", repos)
}
