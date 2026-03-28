package bootstrap

import (
	"fmt"

	"gameclustering.com/internal/core"
)

type EventManager struct {
	App *AppManager
}

func (s *EventManager) Create(classId uint32, topic string) (core.Event, error) {
	return nil, nil
}

func (s *EventManager) VerifyTicket(ticket string) (core.OnSession, error) {
	session, err := s.App.auth.ValidateTicket(ticket)
	if err != nil {
		return session, err
	}
	if session.AccessControl < core.ADMIN_ACCESS_CONTROL {
		return session, fmt.Errorf("admin access control required %d", session.AccessControl)
	}
	return session, nil
}

func (s *EventManager) OnEvent(e core.Event) {
	core.AppLog.Debug().Msgf("event %v", e)
}

func (s *EventManager) OnError(e core.Event, err error) {

}
