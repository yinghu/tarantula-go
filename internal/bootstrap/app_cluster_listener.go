package bootstrap

import (
	"time"

	"gameclustering.com/internal/core"
)

type KVUpdate struct {
	Key      string `json:"Key"`
	Value    string `json:"value"`
	core.Opt `json:"Opt"`
}

func (s *AppManager) MemberJoined(joined core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member joined %v\n", joined)
	go s.sendToApp("presence", "join", joined)
	go s.sendToApp("asset", "join", joined)
	go s.sendToApp("profile", "join", joined)
	go s.sendToApp("inventory", "join", joined)
	go s.sendToApp("tournament", "join", joined)
}
func (s *AppManager) MemberLeft(left core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member left %v\n", left)
	go s.sendToApp("presence", "left", left)
	go s.sendToApp("asset", "left", left)
	go s.sendToApp("profile", "left", left)
	go s.sendToApp("inventory", "left", left)
	go s.sendToApp("tournament", "left", left)
}
func (s *AppManager) Updated(key string, value string, opt core.Opt) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Key updated %s %s %v\n", key, value, opt)
	go s.updateToApp("presence", "update", KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("asset", "update", KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("profile", "update", KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("inventory", "update", KVUpdate{Key: key, Value: value, Opt: opt})
	go s.updateToApp("tournament", "update", KVUpdate{Key: key, Value: value, Opt: opt})
}

func (s *AppManager) updateToApp(app string, cmd string, update KVUpdate) {
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
