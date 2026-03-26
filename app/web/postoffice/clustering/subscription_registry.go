package clustering

import (
	"gameclustering.com/internal/core"
)

type SubscriptionRegistry struct {
	Index string //nodeId | tag | topic | (nodeId:tag:topic)
	Subs  map[string]core.Subscription

	Ends map[core.TopicKey]map[string]core.Subscription
}

func (s *SubscriptionRegistry) add(sub core.Subscription) {
	subs, exists := s.Ends[sub.TopicKey()]
	if !exists {
		subs = make(map[string]core.Subscription)
		s.Ends[sub.TopicKey()] = subs
	}
	_, exists = subs[sub.Key()]
	if exists {
		return
	}
	subs[sub.Key()] = sub
}

func (s *SubscriptionRegistry) del(sub core.Subscription) {
	subs, exists := s.Ends[sub.TopicKey()]
	if !exists {
		return
	}
	delete(subs, sub.Key())
	if len(subs) > 0 {
		return
	}
	delete(s.Ends, sub.TopicKey())
}

func (s *SubscriptionRegistry) size() int {
	return len(s.Ends)
}

func (s *SubscriptionRegistry) topic() []core.Subscription {
	sub := make([]core.Subscription, 0)
	for k := range s.Ends {
		sub = append(sub, core.Subscription{Topic: k.Topic, Endpoint: k.Endpoint})
	}
	return sub
}
