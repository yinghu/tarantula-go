package auth

import (
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
)

type Login struct {
	Name        string `json:"login"`
	Hash        string `json:"password"`
	ReferenceId int32  `json:"referenceId"`
	SystemId    int64

	event.EventObj //Event default
	persistence.PersistentableObj // Persistentable default
}
