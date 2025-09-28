package core

type Authenticator interface {
	Password
	Token
	Ticket
}

type OnSession struct {
	Successful    bool   `json:"successful"`
	ErrorCode     int    `json:"errorCode"`
	Message       string `json:"message"`
	SystemId      int64  `json:"systemId"`
	Stub          int32  `json:"-"`
	Token         string `json:"token"`
	Ticket        string `json:"ticket"`
	Home          string `json:"home"`
	AccessControl int32  `json:"-"`
}
