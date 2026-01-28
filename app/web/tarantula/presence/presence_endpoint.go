package presence

import "net/http"

type PresenceEndpoint struct {
	
}

func (p *PresenceEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	w.Write([]byte("presence node"))
}
