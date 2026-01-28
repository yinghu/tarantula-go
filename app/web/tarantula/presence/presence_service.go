package presence

import (
	"net/http"

	"gameclustering.com/internal/core"
)

func Start(viewer core.ClusterViewer) {
	http.Handle("/tarantula/presence", &PresenceEndpoint{viewer})
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		core.AppLog.Printf("failed to start service %s\n", "presence")
	}
}
