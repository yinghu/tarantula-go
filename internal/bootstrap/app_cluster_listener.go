package bootstrap

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

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
	core.AppLog.Printf("Key updated %s %s %v\n", key, value, opt)
	go s.updateToApp("presence", "update", KVUpdate{Key: key, Value: value, Opt: opt})
}

func (s *AppManager) updateToApp(app string, cmd string, update KVUpdate) {
	tick, err := s.AppAuth.CreateTicket(1, 1, SUDO_ACCESS_CONTROL)
	if err != nil {
		return
	}
	data, err := json.Marshal(update)
	if err != nil {
		return
	}
	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", "http://"+app+":8080/"+app+"/clusteradmin/"+cmd, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		core.AppLog.Printf("Error %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		core.AppLog.Printf("Resp Error %s\n", err.Error())
	}
	core.AppLog.Printf("Response code : %d %s\n", resp.StatusCode, string(r))
}

func (s *AppManager) sendToApp(app string, cmd string, node core.Node) {
	tick, err := s.AppAuth.CreateTicket(1, 1, SUDO_ACCESS_CONTROL)
	if err != nil {
		return
	}
	data, err := json.Marshal(node)
	if err != nil {
		return
	}
	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", "http://"+app+":8080/"+app+"/clusteradmin/"+cmd, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		core.AppLog.Printf("Error %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		core.AppLog.Printf("Resp Error %s\n", err.Error())
	}
	core.AppLog.Printf("Response code : %d %s\n", resp.StatusCode, string(r))
}
