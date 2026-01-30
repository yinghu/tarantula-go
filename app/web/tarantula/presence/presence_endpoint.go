package presence

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
)

type PresenceEndpoint struct {
	core.ClusterViewer
}

func (p *PresenceEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rq := make(chan []core.Node, 1)
	p.HashRing(core.RingRequest{Async: rq})
	n := <-rq
	data, err := json.Marshal(n)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
