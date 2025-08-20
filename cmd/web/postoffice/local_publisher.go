package main

import "gameclustering.com/internal/event"

type LocalPublisher struct {
	*PostofficeService
}

func (s *LocalPublisher) Publish(e event.Event, ticket string) error {
	s.OnEvent(e)
	return nil
}

func (s *LocalPublisher) Close() error {
	return nil
}

func (s *LocalPublisher) Connect() error {
	return nil
}
