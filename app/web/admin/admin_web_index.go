package main

import (
	"io"
	"net/http"
	"os"
	"strings"

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
	fn := r.PathValue("name")
	dest, err := os.OpenFile(s.contentDir+"/site/"+fn, os.O_RDONLY, 0644)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	if strings.HasSuffix(fn, ".js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(fn, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(fn, ".json") {
		w.Header().Set("Content-Type", "application/json")
	} else if strings.HasSuffix(fn, ".ico") {
		w.Header().Set("Content-Type", "image/x-icon")
	} else {
		w.Header().Set("Content-Type", "text/html")
	}
	w.WriteHeader(http.StatusOK)
	defer dest.Close()
	io.Copy(w, dest)
}
