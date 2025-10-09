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
	un := r.PathValue("name")
	ctype := r.Header.Get("Content-type")
	pdir := s.assetDir + "/" + strconv.Itoa(int(rs.SystemId))
	os.MkdirAll(pdir, 0755)
	fid := uuid.New()
	fn := fid.String()+"."+strings.Split(ctype, "/")[1]
	dest, err := os.OpenFile(pdir+"/"+fn, os.O_CREATE|os.O_WRONLY, 0644)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error(), ErrorCode: 1}
		w.Write(util.ToJson(session))
		return
	}
	defer dest.Close()
	rt, err := io.Copy(dest, r.Body)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error(), ErrorCode: 2}
		w.Write(util.ToJson(session))
		return
	}
	err = s.saveAssetIndex(AssetIndex{systemId: rs.SystemId, name: un, fileIndex: fn, uploadTime: time.Now()})
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error(), ErrorCode: 3}
		w.Write(util.ToJson(session))
		return
	}
	session := core.OnSession{Successful: true, Message: "upload [" + strconv.Itoa(int(rt)) + "]"}
	w.Write(util.ToJson(session))
}
