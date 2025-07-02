package bootstrap

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type AppItemAdmin struct {
	TarantulaService
}

func (s *AppItemAdmin) AccessControl() int32 {
	return ADMIN_ACCESS_CONTROL
}

func (s *AppItemAdmin) Request(rs core.OnSession, w http.ResponseWriter, r *http.Request) {
	defer func() {
		w.WriteHeader(http.StatusOK)
		session := core.OnSession{Successful: true, Message: "app item admin [" + s.Cluster().Group() + "]"}
		w.Write(util.ToJson(session))
		r.Body.Close()
	}()
	cmd := r.PathValue("cmd")
	core.AppLog.Printf("command %s\n", cmd)
	var e item.Enum
	json.NewDecoder(r.Body).Decode(&e)
	s.ItemListener().OnEnum(e)
	s.ItemListener().OnCategory(item.Category{})
	s.ItemListener().OnConfiguration(item.Configuration{})
}
