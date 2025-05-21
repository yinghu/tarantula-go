package main

import (
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


func errorMessage(msg string, code int) []byte {
	m := OnSession{Message: msg, ErrorCode: code}
	return util.ToJson(m)
}
func successMessage(msg string) []byte {
	m := OnSession{Message: msg, Successful: true}
	return util.ToJson(m)
}
