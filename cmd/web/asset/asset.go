package main

import (
	"time"

	"gameclustering.com/internal/bootstrap"
)

type AssetIndex struct {
	systemId   int64
	name       string
	fileIndex  string
	uploadTime time.Time
}

func main() {

	bootstrap.AppBootstrap(&AssetService{})

}
