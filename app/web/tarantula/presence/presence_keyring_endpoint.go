package presence

import (
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
)

type PresenceKeyRingEndpoint struct {
	core.ClusterService
}

func (p *PresenceKeyRingEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	k := r.PathValue("key")
	rq := make(chan []core.Node, 1)
	defer close(rq)
	p.KeyRing(core.RingRequest{Token: p.RingToken([]byte(k)), Async: rq})
	n := <-rq
	data, err := json.Marshal(n)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
