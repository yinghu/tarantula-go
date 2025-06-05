package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AssetDownload struct {
	*AssetService
}

func (s *AssetDownload) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *AssetDownload) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	session := core.OnSession{Successful: true, Message: "download"}
	w.Write(util.ToJson(session))
}
