package clustering

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
)

func TestSubscriptionRegistry(t *testing.T) {
	sub1 := core.Subscription{NodeId: "p01", Tag: "presence", Topic: "message", Endpoint: "192.168.1.11:7001"}
	sub2 := core.Subscription{NodeId: "p01", Tag: "presence", Topic: "message", Endpoint: "192.168.1.11:7001"}
	sub3 := core.Subscription{NodeId: "p01", Tag: "tournament", Topic: "score", Endpoint: "192.168.1.11:7001"}
	reg := SubscriptionRegistry{Ends: make(map[core.TopicKey]map[string]core.Subscription)}
	reg.add(sub1)
	reg.add(sub2)
	reg.add(sub3)
	//fmt.Printf("size of subs %v\n",reg.Ends[sub1.Endpoint])
	fmt.Printf("size of reg %d \n", reg.size())
	reg.del(sub3)
	reg.del(sub1)
	reg.del(sub2)
	fmt.Printf("size of reg %d \n", reg.size())

	//fmt.Printf("size of subs %v\n",reg.Ends[sub1.Endpoint])
}
