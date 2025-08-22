package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type InventoryGranter struct {
	*InventoryService
}

func (s *InventoryGranter) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *InventoryGranter) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Write(util.ToJson(core.OnSession{Successful: true, Message: "granted"}))
}
