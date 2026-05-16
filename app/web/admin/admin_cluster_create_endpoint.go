package main

import (
	"bytes"
	"io"
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/encoding/protojson"
)

type AdminClusterCreate struct {
	*AdminService
}

func (s *AdminClusterCreate) AccessControl() int32 {
	return core.PROTECTED_ACCESS_CONTROL
}

func (s *AdminClusterCreate) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	mf := persistence.NewVMObjectFactory()
	var me protocol.VMObject
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r.Body)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	err = protojson.Unmarshal(buf.Bytes(), &me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	kv, err := mf.FromVMObject(&me)
	if err != nil {
		w.Write(util.ToJson(core.OnSession{Successful: false, Message: err.Error()}))
		return
	}
	kv.Key.Array = s.ToBytes(100)
	tb := persistence.NewTaskBuilder(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "register"})

	vb := tb.Validator(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "validator"})
	vb.Transaction().Meta(&protocol.Meta{Name: "check"}).Object(kv).Build()

	jb := tb.Job(&protocol.Meta{NodeId: s.NodeId(), Tag: s.Context(), Name: "job"})
	jb.Transaction().Meta(&protocol.Meta{Name: "create"}).Object(kv).Build()
	jb.Build()
	rp, err := s.Cluster().Issue(tb.Build())
	if err != nil {
		core.AppLog.Debug().Msgf("TASK ERR %s", err.Error())
		return
	}
	core.AppLog.Debug().Msgf("TASK %v", rp)
	w.Write(util.ToJson(core.OnSession{Successful: true, Message: ""}))
}
