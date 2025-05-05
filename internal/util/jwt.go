package util

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"strings"
)

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

type JwtPayload struct {
	Aud string `json:"aud"`
	Exp int64 `json:"exp"`
}

type Jwt struct {
	Mac hash.Hash
	Alg string
	Ksz int16
}

type JwtProcess func(*JwtHeader, *JwtPayload) error

func (j *Jwt) HMac() {
	key := make([]byte,j.Ksz)
	rand.Read(key)
	j.Mac = hmac.New(sha256.New, key)
}
func (j *Jwt) HMacFromKey(key []byte) {
	
	
	j.Mac = hmac.New(sha256.New, key)
}

func (j *Jwt) Token(jp JwtProcess) (string, error) {
	h := JwtHeader{Alg: j.Alg, Typ: "JWT"}
	p := JwtPayload{}
	err := jp(&h, &p)
	if err !=nil {
		return "",err
	}
	var sb strings.Builder
	hd, _ := json.Marshal(h)
	sb.WriteString(base64.URLEncoding.EncodeToString(hd))
	sb.WriteString(".")
	pd, _ := json.Marshal(p)
	sb.WriteString(base64.URLEncoding.EncodeToString(pd))
	j.Mac.Reset()
	j.Mac.Write([]byte(sb.String()))
	sum := j.Mac.Sum(nil)
	sb.WriteString(".")
	sb.WriteString(base64.URLEncoding.EncodeToString(sum))
	return sb.String(), nil
}

func (j *Jwt) Verify(token string, jp JwtProcess) error {
	parts := strings.Split(token, ".")
	var sb strings.Builder
	sb.WriteString(parts[0])
	sb.WriteString(".")
	sb.WriteString(parts[1])
	j.Mac.Reset()
	j.Mac.Write([]byte(sb.String()))
	sum := j.Mac.Sum(nil)
	if base64.URLEncoding.EncodeToString(sum) != parts[2] {
		return errors.New("bad signature")
	}

	h, _ := base64.URLEncoding.DecodeString(parts[0])
	p, _ := base64.URLEncoding.DecodeString(parts[1])

	th := JwtHeader{}
	json.Unmarshal(h, &th)
	tp := JwtPayload{}
	json.Unmarshal(p, &tp)
	err := jp(&th, &tp)
	if err != nil{
		return err
	}
	//fmt.Println(th)
	//fmt.Println(tp)
	//fmt.Println(parts[2])
	return nil
}
