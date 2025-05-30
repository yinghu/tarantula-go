package core

type Authenticator interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password string, hash string) error
	CreateToken(systemId int64, stub int64) (string, error)
	ValidateToken(token string) error
}

type OnSession struct {
	Successful bool   `json:"successful"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	SystemId   int64  `json:"systemId"`
	Stub       int64  `json:"stub"`
	Token      string `json:"token"`
	Home       string `json:"home"`
}
