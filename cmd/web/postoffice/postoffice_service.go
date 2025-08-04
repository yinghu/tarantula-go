package main

import (
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
)

type PostofficeService struct {
	bootstrap.AppManager
	Ds core.DataStore
}

func (s *PostofficeService) Config() string {
	return "/etc/tarantula/postoffice-conf.json"
}

func (s *PostofficeService) Start(env conf.Env, c core.Cluster) error {
	s.AppManager.Start(env, c)
	path := env.LocalDir + "/store"
	ds := persistence.Cache{InMemory: env.Bdg.InMemory, Path: path, Seq: s.Sequence()}
	err := ds.Open()
	if err != nil {
		return err
	}
	s.Ds = &ds
	core.AppLog.Printf("Postoffice service started %s %s\n", env.HttpBinding, env.LocalDir)
	http.Handle("/postoffice/subscribe/{key}", bootstrap.Logging(&PostofficeSubscriber{PostofficeService: s}))
	http.Handle("/postoffice/unsubscribe/{key}", bootstrap.Logging(&PostofficeUnSubscriber{PostofficeService: s}))
	http.Handle("/postoffice/publish/{key}", bootstrap.Logging(&PostofficePublisher{PostofficeService: s}))
	return nil
}

func (s *PostofficeService) Create(classId int, ticket string) (event.Event, error) {
	me := event.CreateEvent(classId, s)
	if me == nil {
		return nil, fmt.Errorf("event ( %d ) not supported", classId)
	}
	return me, nil
}

func (s *PostofficeService) OnError(e error) {
	core.AppLog.Printf("On event error %s\n", e.Error())
}

func (s *PostofficeService) OnEvent(e event.Event) {
	v, ok := e.(*event.MessageEvent)
	if ok {
		core.AppLog.Printf("On event %s, %s\n", v.Message, v.Title)
		s.PostJsonSync(fmt.Sprintf("%s%d", "http://tournament:8080/tournament/clusteradmin/event/", v.ClassId()), v)
	}
}
