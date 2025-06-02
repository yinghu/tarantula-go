package main

import (
	"gameclustering.com/internal/bootstrap"
)

const (
	DB_OP_ERR_CODE     int    = 500100
	WRONG_PASS_CODE    int    = 400100
	WRONG_PASS_MSG     string = "wrong user/password"
	INVALID_TOKEN_CODE int    = 400101
	INVALID_TOKEN_MSG  string = "invalid token"
)

func main() {

	bootstrap.AppBootstrap(&PresenceService{})

}
