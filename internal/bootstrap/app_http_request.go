package bootstrap

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gameclustering.com/internal/core"
)

func (s *AppManager) PostJson(url string, payload any, ch chan core.OnSession) {
	tick, err := s.AppAuth.CreateTicket(1, 1, SUDO_ACCESS_CONTROL)
	if err != nil {
		ch <- core.OnSession{ErrorCode: INVALID_TICKET_CODE, Message: err.Error()}
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		ch <- core.OnSession{ErrorCode: BAD_REQUEST_CODE, Message: err.Error()}
		return
	}
	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		ch <- core.OnSession{ErrorCode: INVALID_JSON_CODE, Message: err.Error()}
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		ch <- core.OnSession{ErrorCode: BAD_REQUEST_CODE, Message: err.Error()}
		return
	}
	defer resp.Body.Close()
	var rt core.OnSession
	err = json.NewDecoder(resp.Body).Decode(&rt)
	if err != nil {
		ch <- core.OnSession{ErrorCode: INVALID_JSON_CODE, Message: err.Error()}
		core.AppLog.Printf("Resp Error %s\n", err.Error())
		return
	}
	ch <- rt
	core.AppLog.Printf("Response code : %d %v\n", resp.StatusCode, rt)
}
