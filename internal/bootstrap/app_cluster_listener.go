package bootstrap

import (
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (s *AppManager) MemberJoined(joined core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member joined %v\n", joined)
	for i := range s.ManagedApps {
		go s.sendToApp(s.ManagedApps[i], "join", joined)
	}
}
func (s *AppManager) MemberLeft(left core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member left %v\n", left)
	for i := range s.ManagedApps {
		go s.sendToApp(s.ManagedApps[i], "left", left)
	}
}
func (s *AppManager) KVUpdated(key string, value string, opt core.Opt) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Key updated %s %s %v\n", key, value, opt)
	go s.updateToApp("presence", "update", item.KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("asset", "update", item.KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("profile", "update", item.KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("inventory", "update", item.KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("tournament", "update", item.KVUpdate{Key: key, Value: value, Opt: opt})
}

func (s *AppManager) updateToApp(app string, cmd string, update item.KVUpdate) {
	for i := range 5 {
		ret := s.PostJsonSync("http://"+app+":8080/"+app+"/clusteradmin/"+cmd, update)
		if ret.ErrorCode == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
	}
}

func (s *AppManager) sendToApp(app string, cmd string, node core.Node) {
	for i := range 5 {
		ret := s.PostJsonSync("http://"+app+":8080/"+app+"/clusteradmin/"+cmd, node)
		if ret.ErrorCode == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
	}
}
