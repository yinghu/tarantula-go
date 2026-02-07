package presence

import (
	"net/http"

	"gameclustering.com/internal/core"
)

type PresenceKeyValueEndpoint struct {
	core.ClusterService
}

func (p *PresenceKeyValueEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	k := r.PathValue("key")
	rq := make(chan core.Chunk, 3)
	defer close(rq)
	p.Get(core.GetRequest{Key: []byte(k), Async: rq})
	for c := range rq {
		if len(c.Data) > 0 {
			w.Write(c.Data)
		}
		if !c.Remaining {
			break
		}
	}
}
