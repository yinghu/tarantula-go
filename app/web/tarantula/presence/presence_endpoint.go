package presence

import (
	"net/http"

	"gameclustering.com/internal/core"
)

type PresenceEndpoint struct {
	core.ClusterViewer
}

func (p *PresenceEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	p.List()
	w.Write([]byte("presence node"))
}
