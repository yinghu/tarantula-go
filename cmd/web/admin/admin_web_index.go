package main

import (
	"io"
	"net/http"
	"os"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type AdminWebIndex struct {
	*AdminService
}

func (s *AdminWebIndex) AccessControl() int32 {
	return bootstrap.PUBLIC_ACCESS_CONTROL
}
func (s *AdminWebIndex) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fn := r.PathValue("name");
	dest, err := os.OpenFile("site/"+fn, os.O_RDONLY, 0644)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	defer dest.Close()
	io.Copy(w, dest)
}
