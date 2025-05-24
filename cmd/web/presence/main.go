package main

import (

	"gameclustering.com/internal/bootstrap"
	
)





func main() {
	
	bootstrap.WebBootstrap(&PresenceService{})
	
}
