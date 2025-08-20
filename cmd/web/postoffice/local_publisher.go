package main

import "gameclustering.com/internal/event"

type LocalPublisher struct {
	*PostofficeService
}

func (s *LocalPublisher) Publish(e event.Event, ticket string) {
	s.OnEvent(e)
}

func (s *LocalPublisher) Close() error{
	return nil
}