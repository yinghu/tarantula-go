package bootstrap

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"gameclustering.com/internal/core"
)

func (s *AppManager) MemberJoined(joined core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member joined %v\n", joined)
	tick, err := s.AppAuth.CreateTicket(1, 1, SUDO_ACCESS_CONTROL)
	if err != nil {
		return
	}
	data, err := json.Marshal(joined)
	if err != nil {
		return
	}
	go sendToApp(tick, "presence", "join", data)
	go sendToApp(tick, "asset", "join", data)
	go sendToApp(tick, "profile", "join", data)
	go sendToApp(tick, "inventory", "join", data)
	go sendToApp(tick, "tournament", "join", data)
}
func (s *AppManager) MemberLeft(left core.Node) {
	if s.standalone {
		return
	}
	core.AppLog.Printf("Member left %v\n", left)
	tick, err := s.AppAuth.CreateTicket(1, 1, SUDO_ACCESS_CONTROL)
	if err != nil {
		return
	}
	data, err := json.Marshal(left)
	if err != nil {
		return
	}
	go sendToApp(tick, "presence", "left", data)
	go sendToApp(tick, "asset", "left", data)
	go sendToApp(tick, "profile", "left", data)
	go sendToApp(tick, "inventory", "left", data)
	go sendToApp(tick, "tournament", "left", data)
}
func (s *AppManager) Updated(key string, value string) {
	core.AppLog.Printf("Key updated %s %s\n", key, value)
}

func sendToApp(ticket string, app string, cmd string, data []byte) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://web/"+app+"/admin/"+cmd, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ticket)
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
