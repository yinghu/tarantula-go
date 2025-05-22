package core

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

type JwtPayload struct {
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
}

type JwtProcess func(*JwtHeader, *JwtPayload) error

type Jwt interface {
	Token(jp JwtProcess) (string, error)
	Verify(token string, jp JwtProcess) error
}
