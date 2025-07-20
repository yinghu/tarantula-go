package main

import (
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type CategoryPublisher struct {
	*AdminService
}

func (s *CategoryPublisher) AccessControl() int32 {
	return bootstrap.ADMIN_ACCESS_CONTROL
}
func (s *CategoryPublisher) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.WriteHeader(http.StatusOK)
	cid, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	os.Chdir(s.publishDir);
	cmd := exec.Command("git","status")
	output, err := cmd.Output()
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	core.AppLog.Printf("Publish category :%d\n", cid)
	session := core.OnSession{Successful: true, Message: string(output)}
	w.Write(util.ToJson(session))
}
