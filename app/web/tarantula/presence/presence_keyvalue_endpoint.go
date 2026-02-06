package presence

import (
	"net/http"

	"gameclustering.com/internal/core"
)

type PresenceKeyValueEndpoint struct {
	core.ClusterViewer
}

func (p *PresenceKeyValueEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	k := r.PathValue("key")
	rq := make(chan core.Chunk,3)
	defer close(rq)
	p.FindValue(core.ValueRequest{Key: []byte(k), Async: rq})
	for c := range rq {
		if len(c.Data) > 0 {
			w.Write(c.Data)
		}
		if !c.Remaining {
			break
		}
	}
}
