package main

import (
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminClusterUpload struct {
	*AdminService
}

func (s *AdminClusterUpload) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *AdminClusterUpload) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	provider := r.PathValue("provider")
	
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: provider}))
}
