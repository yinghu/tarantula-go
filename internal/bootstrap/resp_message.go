package bootstrap

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

func ErrorMessage(msg string, code int) []byte {
	m := core.OnSession{Message: msg, ErrorCode: code}
	return util.ToJson(m)
}
func SuccessMessage(msg string) []byte {
	m := core.OnSession{Message: msg, Successful: true}
	return util.ToJson(m)
}
