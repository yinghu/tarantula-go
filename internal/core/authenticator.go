package core

type Authenticator interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password string, hash string) error
	CreateToken(systemId int64, stub int32, accessControl int32) (string, error)
	ValidateToken(token string) (OnSession, error)
}

type OnSession struct {
	Successful    bool   `json:"successful"`
	ErrorCode     int    `json:"errorCode"`
	Message       string `json:"message"`
	SystemId      int64  `json:"systemId"`
	Stub          int32  `json:"-"`
	Token         string `json:"token"`
	Home          string `json:"home"`
	AccessControl int32  `json:"-"`
}
