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
	nodes := p.List()
	data, err := json.Marshal(nodes)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
