package main

import "gameclustering.com/internal/core"

type LocalPublisher struct {
	*PostofficeService
}

func (s *LocalPublisher) Publish(e core.Event, ticket string) error {
	s.Event().OnEvent(e)
	return nil
}

func (s *LocalPublisher) Close() error {
	return nil
}

func (s *LocalPublisher) Connect() error {
	return nil
}
