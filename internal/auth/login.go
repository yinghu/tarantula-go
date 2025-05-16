package auth

import (
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
)

const (
	DB_OP_ERR_CODE int = 500100

	WRONG_PASS_CODE    int    = 400100
	WRONG_PASS_MSG     string = "wrong user/password"
	INVALID_TOKEN_CODE int    = 400101
	INVALID_TOKEN_MSG  string = "invalid token"
)

type OnSession struct {
	Successful bool   `json:"successful"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	SystemId   int64  `json:"systemId"`
	Stub       int64  `json:"stub"`
	Token      string `json:"token"`
	Home       string `json:"home"`
}

type Login struct {
	Name        string `json:"login"`
	Hash        string `json:"password"`
	ReferenceId int32  `json:"referenceId"`
	SystemId    int64

	event.EventObj                //Event default
	persistence.PersistentableObj // Persistentable default
}

func errorMessage(msg string, code int) []byte {
	m := OnSession{Message: msg, ErrorCode: code}
	return util.ToJson(m)
}
func successMessage(msg string) []byte {
	m := OnSession{Message: msg, Successful: true}
	return util.ToJson(m)
}
