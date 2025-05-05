package util

import (
	"errors"
	"testing"
	"time"
	"crypto/rand"
)

func onToken(h *JwtHeader, p *JwtPayload) error {
	h.Kid = "kid"
	p.Aud = "player"
	exp := time.Now().Add(time.Hour*24).UTC()
	//fmt.Println(exp)
	p.Exp = exp.UnixMilli()
	return nil
}

func onVerify(h *JwtHeader, p *JwtPayload) error {
	//fmt.Println(h.Kid)
	//fmt.Println(p.Aud)
	t := time.UnixMilli(p.Exp).UTC()
	if t.Before(time.Now().UTC()) {
		return errors.New("token expired")
	}
	//fmt.Println(t)
	return nil
}

func TestSum(t *testing.T) {
	jwt := Jwt{Alg: "SHS256"}
	key := make([]byte,32)
	rand.Read(key)
	jwt.HMacFromKey(key)
	token, _ := jwt.Token(onToken)
	time.Sleep(1000*time.Millisecond)
	err := jwt.Verify(token, onVerify)
	if err != nil{
		t.Errorf("failed %s\n",err.Error())
	}
}
