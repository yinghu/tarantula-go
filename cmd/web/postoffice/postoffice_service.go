package main

import (
	"fmt"
	"net/http"
	"sync"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/conf"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
)

type TopicMap struct {
	sync.RWMutex
	topics map[int32]event.SubscriptionEvent
}

type PostofficeService struct {
	bootstrap.AppManager
	Ds core.DataStore
	TopicMap
}

func (s *PostofficeService) Config() string {
	return "/etc/tarantula/postoffice-conf.json"
}

func (s *PostofficeService) Start(env conf.Env, c core.Cluster) error {
	s.AppManager.Start(env, c)
	s.createSchema()
	path := env.LocalDir + "/store"
	ds := persistence.Cache{InMemory: env.Bdg.InMemory, Path: path, Seq: s.Sequence()}
	err := ds.Open()
	if err != nil {
		return err
	}
	s.Ds = &ds
	s.topics = make(map[int32]event.SubscriptionEvent)
	s.loadTopics()
	core.AppLog.Printf("Postoffice service started %s %s\n", env.HttpBinding, env.LocalDir)
	http.Handle("/postoffice/subscribe", bootstrap.Logging(&PostofficeSubscriber{PostofficeService: s}))
	http.Handle("/postoffice/unsubscribe", bootstrap.Logging(&PostofficeUnSubscriber{PostofficeService: s}))
	http.Handle("/postoffice/publish/{topic}/{cid}", bootstrap.Logging(&PostofficePublisher{PostofficeService: s}))
	return nil
}

func (s *PostofficeService) Create(classId int, topic string) (event.Event, error) {
	me := event.CreateEvent(classId, s)
	me.Topic(topic)
	if me == nil {
		return nil, fmt.Errorf("event ( %d ) not supported", classId)
	}
	return me, nil
}

func (s *PostofficeService) OnError(e error) {
	core.AppLog.Printf("On event error %s\n", e.Error())
}

func (s *PostofficeService) OnEvent(e event.Event) {
	se, isSe := e.(*event.SubscriptionEvent)
	if isSe {
		core.AppLog.Printf("On event %d %s, %s, %s %d\n", se.Id, se.App, se.Name, se.OnTopic(), se.ClassId())
		s.RWMutex.Lock()
		defer s.RWMutex.Unlock()
		s.topics[se.Id] = *se
		return
	}
	s.RLock()
	defer s.RUnlock()
	apps := make([]string, 0)
	for i := range s.topics {
		if s.topics[i].Name != e.OnTopic() {
			continue
		}
		apps = append(apps, s.topics[i].App)
	}
	go func() {
		for x := range apps {
			url := fmt.Sprintf("%s%s%s%s%s%d", "http://", apps[x], ":8080/", apps[x], "/clusteradmin/event/", e.ClassId())
			core.AppLog.Printf("Pushlish to %s\n", url)
			s.PostJsonSync(url, e)
		}
	}()

}

func (s *PostofficeService) Publish(e event.Event) {
	ticket, err := s.AppAuth.CreateTicket(0, 0, bootstrap.ADMIN_ACCESS_CONTROL)
	if err != nil {
		core.AppLog.Printf("Ticket error %s\n", err.Error())
		return
	}
	view := s.Cluster().View()
	for i := range view {
		v := view[i]
		core.AppLog.Printf("Sending to : %s,%s,%s,%s\n", v.Name, v.TcpEndpoint, s.Cluster().Local().Name, e.OnETag()) // no prefix
		if v.Name == s.Cluster().Local().Name {
			s.OnEvent(e)
			continue
		}
		pub := event.SocketPublisher{Remote: v.TcpEndpoint}
		pub.Publish(e, ticket)
	}
}
