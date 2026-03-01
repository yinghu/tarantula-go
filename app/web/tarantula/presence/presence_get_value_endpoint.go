package presence

import (
	"net/http"
	"strconv"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type PresenceGetValueEndpoint struct {
	core.ClusterService
}

func (p *PresenceGetValueEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	k := r.PathValue("key")
	f, _ := strconv.ParseInt(r.PathValue("factoryId"), 10, 32)
	c, _ := strconv.ParseInt(r.PathValue("classId"), 10, 32)

	rq := make(chan core.Chunk, 3)
	defer close(rq)
	req := core.GetRequest{Key: []byte(k), Async: rq}
	req.FactoryId = int32(f)
	req.ClassId = int32(c)
	core.AppLog.Debug().Msg("Calling get api")
	p.Get(req)
	for c := range rq {
		if len(c.Data) > 0 {
			w.Write(util.ToJson(c))
		}
		if !c.Remaining {
			break
		}
	}
}
