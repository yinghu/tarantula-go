package main

import (
	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/event"
)

type Login struct {
	Id            int32            `json:"-"`
	Name          string           `json:"login"`
	Hash          string           `json:"password"`
	ReferenceId   int32            `json:"referenceId"`
	SystemId      int64            `json:"systemId:string"`
	AccessControl int32            `json:"accessControl,string"`
	Cc            chan event.Chunk `json:"-"`
}

func main() {

	bootstrap.AppBootstrap(&PresenceService{})

}
