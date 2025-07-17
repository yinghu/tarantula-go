package main

import (
	"io"
	"net/http"
	"os"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
	"github.com/google/uuid"
)

type FileSaver struct {
	*AdminService
}

func (s *FileSaver) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *FileSaver) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fid := uuid.New()
	fn := fid.String() + ".json"
	dest, err := os.OpenFile(s.assetDir+"/"+fn, os.O_CREATE|os.O_WRONLY, 0644)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error(), ErrorCode: 1}
		w.Write(util.ToJson(session))
		return
	}
	defer dest.Close()
	_, err = io.Copy(dest, r.Body)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error(), ErrorCode: 2}
		w.Write(util.ToJson(session))
		return
	}
	session := core.OnSession{Successful: true, Message: fn}
	w.Write(util.ToJson(session))
}
