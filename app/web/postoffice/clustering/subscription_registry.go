package clustering

import (
	"gameclustering.com/internal/core"
)

type SubscriptionRegistry struct {
	
	topicEnds map[core.TopicKey]map[string]core.Subscription
}

func (s *SubscriptionRegistry) add(sub core.Subscription) {
	subs, exists := s.topicEnds[sub.TopicKey()]
	if !exists {
		subs = make(map[string]core.Subscription)
		s.topicEnds[sub.TopicKey()] = subs
	}
	_, exists = subs[sub.Key()]
	if exists {
		return
	}
	subs[sub.Key()] = sub
}

func (s *SubscriptionRegistry) del(sub core.Subscription) {
	subs, exists := s.topicEnds[sub.TopicKey()]
	if !exists {
		return
	}
	delete(subs, sub.Key())
	if len(subs) > 0 {
		return
	}
	delete(s.topicEnds, sub.TopicKey())
}

func (s *SubscriptionRegistry) size() int {
	return len(s.topicEnds)
}

func (s *SubscriptionRegistry) topic(name string) []core.Subscription {
	sub := make([]core.Subscription, 0)
	for k := range s.topicEnds {
		if name == k.Topic {
			sub = append(sub, core.Subscription{Topic: k.Topic, Endpoint: k.Endpoint})
		}
	}
	return sub
}
