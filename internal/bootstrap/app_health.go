package bootstrap

import (
	"net/http"

	"gameclustering.com/internal/core"
)

type AppHealth struct {
	TarantulaService
}

func (s *AppHealth) AccessControl() int32 {
	return PUBLIC_ACCESS_CONTROL
}

func (s *AppHealth) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
