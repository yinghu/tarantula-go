package clustering

type TopicRequest struct {
	Topic     string
	Listeners []string
}

type SubscriptionRegistry struct {
	subs     map[string][]string
	register chan TopicRequest
}

func (s SubscriptionRegistry) Monitor(){
	
}
