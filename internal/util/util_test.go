package util

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
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

func TestTick(t *testing.T) {
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	c := 3
	for range tick.C {
	rp:
		c--
		if c > 0 {
			fmt.Printf("tick %v\n", time.Now())
			goto rp
		}
		c = 3
	}
}
