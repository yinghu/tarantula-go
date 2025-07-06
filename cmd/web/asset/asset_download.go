package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

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
	un := r.PathValue("name")
	aIndex := AssetIndex{systemId: rs.SystemId, name: un}
	err := s.loadAssetIndex(&aIndex)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	pdir := s.assetDir + "/" + strconv.Itoa(int(rs.SystemId))
	dest, err := os.OpenFile(pdir+"/"+aIndex.fileIndex, os.O_RDONLY, 0644)
	if err != nil {
		session := core.OnSession{Successful: false, Message: err.Error()}
		w.Write(util.ToJson(session))
		return
	}
	defer dest.Close()
	io.Copy(w, dest)
}
