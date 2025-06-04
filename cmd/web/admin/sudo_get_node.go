package main

import (
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/cluster"
	"gameclustering.com/internal/core"
)

type AdminGetNode struct {
	*AdminService
}

func (s *AdminGetNode) AccessControl() int32 {
	return bootstrap.SUDO_ACCESS_CONTROL
}
func (s *AdminGetNode) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	group := r.PathValue("group")
	node := r.PathValue("name")
	var value string
	err := s.Cluster().AtomicWithPrefix(group, func(ctx cluster.Ctx) error {
		val, err := ctx.Get(node)
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	if err != nil {
		w.Write(bootstrap.ErrorMessage("node not existed", 400001))
		return
	}
	w.Write([]byte(value))
}
