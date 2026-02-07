package presence

import (
	"net/http"

	"gameclustering.com/internal/core"
)

func Start(viewer core.ClusterService) {
	http.Handle("/tarantula/presence/hashring", &PresenceHashRingEndpoint{viewer})
	http.Handle("/tarantula/presence/keyring/{key}", &PresenceKeyRingEndpoint{viewer})
	http.Handle("/tarantula/presence/keyvalue/{key}", &PresenceKeyValueEndpoint{viewer})
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		core.AppLog.Printf("failed to start service %s\n", "presence")
	}
}
