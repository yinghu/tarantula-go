package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

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

type CChange struct {
	nodeName string
	endpoint string
	started  bool
}

type PostofficeService struct {
	bootstrap.AppManager
	Ds core.DataStore
	TopicMap
	outboundQ chan event.Event
	cchangeQ  []chan CChange

	inboundQ chan event.Event
	topicQ   []chan event.SubscriptionEvent
	ready    sync.WaitGroup
}

func (s *PostofficeService) Config() string {
	return "/etc/tarantula/postoffice-conf.json"
}

func (s *PostofficeService) Start(env conf.Env, c core.Cluster, p event.Pusher) error {
	s.Bsl = s
	s.AppManager.Start(env, c, p)
	s.createSchema()
	s.topics = make(map[int32]event.SubscriptionEvent)
	s.loadTopics()
	s.ready = sync.WaitGroup{}
	s.ready.Add(1)
	s.outboundQ = make(chan event.Event, 10)
	s.cchangeQ = make([]chan CChange, 0)
	ec := make(chan CChange, 1)
	s.cchangeQ = append(s.cchangeQ, ec)
	go s.outboundEvent(ec)

	s.inboundQ = make(chan event.Event, 10)
	s.topicQ = make([]chan event.SubscriptionEvent, 0)
	tc := make(chan event.SubscriptionEvent, 1)
	s.topicQ = append(s.topicQ, tc)
	go s.inboundEvent(tc)

	path := env.LocalDir + "/store"
	ds := persistence.BadgerLocal{InMemory: env.Bdg.InMemory, Path: path, GcEnabled: true}
	err := ds.Open()
	if err != nil {
		return err
	}
	s.Ds = &ds
	s.ready.Done()
	core.AppLog.Printf("Postoffice service started %s %s\n", env.HttpBinding, env.LocalDir)
	http.Handle("/postoffice/subscribe", bootstrap.Logging(&PostofficeSubscriber{PostofficeService: s}))
	http.Handle("/postoffice/unsubscribe", bootstrap.Logging(&PostofficeUnSubscriber{PostofficeService: s}))
	http.Handle("/postoffice/publish/{topic}/{cid}", bootstrap.Logging(&PostofficePublisher{PostofficeService: s}))
	http.Handle("/postoffice/query/{id}", bootstrap.Logging(&PostofficeQueryer{PostofficeService: s}))
	http.Handle("/postoffice/recover/{id}", bootstrap.Logging(&PostofficeRecoverer{PostofficeService: s}))
	return nil
}

func (s *PostofficeService) Shutdown() {
	s.Ds.Close()
	s.AppManager.Shutdown()
	core.AppLog.Println("psotoffice service shutdown")
}

func (s *PostofficeService) Create(classId int, topic string) (event.Event, error) {
	me := event.CreateEvent(classId)
	me.OnListener(s)
	me.OnTopic(topic)
	if me == nil {
		return nil, fmt.Errorf("event ( %d ) not supported", classId)
	}
	return me, nil
}

func (s *PostofficeService) OnError(e event.Event, err error) {
	core.AppLog.Printf("On event error %s\n", err.Error())
}

func (s *PostofficeService) OnEvent(e event.Event) {
	se, isSe := e.(*event.SubscriptionEvent)
	if isSe {
		for i := range s.topicQ {
			s.topicQ[i] <- *se
		}
		return
	}
	s.inboundQ <- e
}
func (s *PostofficeService) LocalStore() core.DataStore {
	return s.Ds
}
func (s *PostofficeService) Publish(e event.Event) {
	s.outboundQ <- e
}

func (s *PostofficeService) NodeStarted(n core.Node) {
	core.AppLog.Printf("node started : %s\n", n.TcpEndpoint)
	for i := range s.cchangeQ {
		s.cchangeQ[i] <- CChange{nodeName: n.Name, endpoint: n.TcpEndpoint, started: true}
	}
}

func (s *PostofficeService) NodeStopped(n core.Node) {
	core.AppLog.Printf("node stopped : %s\n", n.TcpEndpoint)
	for i := range s.cchangeQ {
		s.cchangeQ[i] <- CChange{nodeName: n.Name, started: false}
	}
}

func (s *PostofficeService) onRetry(e event.Event) {
	core.AppLog.Printf("Retrying %v\n", e)
}
func (s *PostofficeService) inboundEvent(t chan event.SubscriptionEvent) {
	topics := make([]event.SubscriptionEvent, 0)
	s.ready.Wait()
	core.AppLog.Printf("inbound queue is ready with pending event size %d\n", len(s.inboundQ))
	s.RLock()
	for _, t := range s.topics {
		topics = append(topics, t)
	}
	s.RUnlock()
	for {
		select {
		case c := <-t:
			topics = append(topics, c)
		case e := <-s.inboundQ:
			core.AppLog.Printf("Inbound event %v\n", e)
			if err := s.Ds.Save(e); err == nil {
				e.OnIndex(s)
			}
			for i := range topics {
				topic := topics[i]
				if topic.Name != e.Topic() {
					continue
				}
				url := fmt.Sprintf("%s%s%s%s%s%d", "http://", topic.App, ":8080/", topic.App, "/clusteradmin/event/", e.ClassId())
				//core.AppLog.Printf("Pushlish to %s\n", url)
				s.PostJsonSync(url, e)
			}
		}
	}
}
func (s *PostofficeService) outboundEvent(c chan CChange) {
	pubs := make(map[string]event.Publisher)
	localListener := LocalEventListener{s}
	for {
		select {
		case e := <-s.outboundQ:
			ticket, err := s.AppAuth.CreateTicket(0, 0, bootstrap.ADMIN_ACCESS_CONTROL)
			if err != nil {
				core.AppLog.Printf("Ticket error %s\n", err.Error())
				continue
			}
			retrying := true
			for _, pub := range pubs {
				e.OnListener(&localListener)
				for i := range 3 {
					if err := pub.Publish(e, ticket); err != nil {
						//break
						core.AppLog.Printf("reconnect to %s retries: %d", err.Error(), i)
						time.Sleep(500 * time.Millisecond)
						pub.Connect()
						continue
					}
					break
				}
				retrying = false
			}
			if retrying {
				s.onRetry(e)
			}
		case c := <-c:
			core.AppLog.Printf("Node Updated : %v\n", c)
			if c.started {
				if c.nodeName == s.Cluster().Local().Name {
					pubs[c.nodeName] = &LocalPublisher{s}
				} else {
					sb := event.SocketPublisher{Remote: c.endpoint}
					for i := range 5 {
						err := sb.Connect()
						if err != nil {
							core.AppLog.Printf("cannot to dial to %s retries: %d", err.Error(), i)
							time.Sleep(1 * time.Second)
							continue
						}
						core.AppLog.Printf("connected %s\n", c.endpoint)
						pubs[c.nodeName] = &sb
						break
					}
				}
			} else {
				pub, rd := pubs[c.nodeName]
				if !rd {
					continue
				}
				pub.Close()
				delete(pubs, c.nodeName)
			}
		}
	}
}
