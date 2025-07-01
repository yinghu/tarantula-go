package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/core"
)

type AppItemAdmin struct {
	TarantulaService
}

func (s *AppItemAdmin) AccessControl() int32 {
	return ADMIN_ACCESS_CONTROL
}

func (s *AppItemAdmin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {

	
}
