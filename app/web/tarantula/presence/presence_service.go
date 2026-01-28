package presence

import (
	"net/http"

	"gameclustering.com/internal/core"
)

func Start() {
	http.Handle("/tarantula/presence", &PresenceEndpoint{})
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		core.AppLog.Printf("failed to start service %s\n", "presence")
	}
}
