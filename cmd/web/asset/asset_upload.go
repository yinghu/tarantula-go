package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"

	"github.com/google/uuid"
)

type AssetUpload struct {
	*AssetService
}

func (s *AssetUpload) AccessControl() int32 {
	return bootstrap.PROTECTED_ACCESS_CONTROL
}

func (s *AssetUpload) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctype := r.Header.Get("Content-type")
	pdir := s.assetDir + "/" + strconv.Itoa(int(rs.SystemId))
	os.MkdirAll(pdir, 0755)
	fid := uuid.New()
	dest, err := os.OpenFile(pdir+"/"+fid.String()+"."+strings.Split(ctype, "/")[1], os.O_CREATE, 0644)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	defer dest.Close()
	rt, err := io.Copy(dest, r.Body)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	err = s.saveAssetIndex(AssetIndex{systemId: rs.SystemId, name: "profile.png", fileIndex: fid.String(), uploadTime: time.Now()})
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	session := core.OnSession{Successful: true, Message: "upload [" + strconv.Itoa(int(rt)) + "]"}
	w.Write(util.ToJson(session))
}
