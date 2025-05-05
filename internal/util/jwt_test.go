package util

import (
	"fmt"
	"testing"
	"time"
)

func onToken(h *JwtHeader, p *JwtPayload) error {
	h.Kid = "kid"
	p.Aud = "player"
	exp := time.Now().Add(time.Hour*24).UTC()
	fmt.Println(exp)
	p.Exp = exp.UnixMilli()
	return nil
}

func onVerify(h *JwtHeader, p *JwtPayload) error {
	fmt.Println(h.Kid)
	fmt.Println(p.Aud)
	t := time.UnixMilli(p.Exp).UTC()
	fmt.Println(t)
	return nil
}

func TestSum(t *testing.T) {
	jwt := Jwt{Alg: "SHS256"}
	jwt.HMac()
	token, _ := jwt.Token(onToken)
	jwt.Verify(token, onVerify)
}
