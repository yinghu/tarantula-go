package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/util"
	"gameclustering.com/internal/core"
)

const (
	DB_OP_ERR_CODE int = 500100

	WRONG_PASS_CODE    int    = 400100
	WRONG_PASS_MSG     string = "wrong user/password"
	INVALID_TOKEN_CODE int    = 400101
	INVALID_TOKEN_MSG  string = "invalid token"
)


func errorMessage(msg string, code int) []byte {
	m := core.OnSession{Message: msg, ErrorCode: code}
	return util.ToJson(m)
}
func successMessage(msg string) []byte {
	m := core.OnSession{Message: msg, Successful: true}
	return util.ToJson(m)
}

func main() {

	bootstrap.AppBootstrap(&PresenceService{})

}
