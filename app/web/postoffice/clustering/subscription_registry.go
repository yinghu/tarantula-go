package clustering

import (
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/protocol"
)

const (
	RECEIVER_START uint32 = 1
	TOPIC_REGISTER uint32 = 2
)

type TopicRequest struct {
	Opt  uint32
	Name string
	Rev  chan chan *protocol.Response
}

type SubscriptionRegistry struct {
	listeners     map[string]chan *protocol.Response
	subscriptions map[string][]core.Subscription
	register      chan core.Subscription
	request       chan TopicRequest
	Messager      chan *protocol.Response
	running       bool
}

func (s SubscriptionRegistry) Register() {
	core.AppLog.Warn().Msg("subscription registry started")
	s.register = make(chan core.Subscription, NODE_EVENT_BUFFER_SIZE)
	s.request = make(chan TopicRequest, NODE_EVENT_BUFFER_SIZE)
	s.Messager = make(chan *protocol.Response, NODE_EVENT_BUFFER_SIZE)
	for s.running {
		select {
		case sub := <-s.register:
			core.AppLog.Debug().Msgf("subscription %v", sub)
			esb, ok := s.subscriptions[sub.Topic]
			if !ok {
				esb = make([]core.Subscription, 0)
			}
			esb = append(esb, sub)
			s.subscriptions[sub.Topic] = esb
		case req := <-s.request:
			core.AppLog.Debug().Msgf("request %v", req)
			switch req.Opt {
			case RECEIVER_START:
				rev := make(chan *protocol.Response, NODE_EVENT_BUFFER_SIZE)
				s.listeners[req.Name] = rev
				req.Rev <- rev
			case TOPIC_REGISTER:

			}
		case msg := <-s.Messager:
			for n, ch := range s.listeners {
				core.AppLog.Debug().Msgf("send message to %s", n)
				ch <- msg
			}
		}
	}
	core.AppLog.Warn().Msg("subscription registry stopped")
}
