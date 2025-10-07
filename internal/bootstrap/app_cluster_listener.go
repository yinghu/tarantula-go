package bootstrap

import (
	"strings"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
)

func (s *AppManager) MemberJoined(joined core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member joined %v\n", joined)
	joined.Local = s.cls.Local().Name == joined.Name
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
	hs := strings.Split(key, ":")
	core.AppLog.Printf("Key updated %s %s %v %d\n", key, value, opt, len(hs))
	if s.standalone {
		return
	}
	n := len(hs)
	if n == 1 {
		for i := range s.ManagedApps {
			go s.updateToApp(s.ManagedApps[i], "update", item.KVUpdate{Key: key, Value: value, Opt: opt})
		}
		return
	}
	//key > {itemId}:{app}
	ns := strings.Split(hs[1], ",")
	for i := range ns {
		go s.updateToApp(ns[i], "schedule", item.KVUpdate{Key: hs[0], Value: value, Opt: opt})
	}
}

func (s *AppManager) updateToApp(app string, cmd string, update item.KVUpdate) {
	for i := range 5 {
		ret := s.PostJsonSync("http://"+app+":8080/"+app+"/clusteradmin/"+cmd+"/0", update)
		if ret.ErrorCode == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
	}
}

func (s *AppManager) sendToApp(app string, cmd string, node core.Node) {
	for i := range 5 {
		ret := s.PostJsonSync("http://"+app+":8080/"+app+"/clusteradmin/"+cmd+"/0", node)
		if ret.ErrorCode == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
	}
}
