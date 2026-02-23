package presence

import (
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type PresenceSetValueEndpoint struct {
	core.ClusterService
}

func (p *PresenceSetValueEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	k := r.PathValue("key")
	v := r.PathValue("value")
	rq := make(chan core.Chunk, 3)
	defer close(rq)
	p.Set(core.SetRequest{Key: []byte(k), Value: []byte(v), Async: rq})
	for c := range rq {
		if len(c.Data) > 0 {
			w.Write(util.ToJson(c))
		}
		if !c.Remaining {
			break
		}
	}
}
