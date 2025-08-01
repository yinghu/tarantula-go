package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
)

type CacheGetter struct {
	*CacheService
}

func (s *CacheGetter) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}

func (s *CacheGetter) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {

}
