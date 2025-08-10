package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

func (s *AppManager) GetJsonAsync(url string, ch chan event.Chunk) {
	if s.standalone {
		ch <- event.Chunk{Remaining: false, Data: util.ToJson(core.OnSession{ErrorCode: STANDALONE_APP, Message: STANDALONE_APP_MSG})}
		return
	}
	tick, err := s.AppAuth.CreateTicket(1, 1, ADMIN_ACCESS_CONTROL)
	if err != nil {
		ch <- event.Chunk{Remaining: false, Data: util.ToJson(core.OnSession{ErrorCode: INVALID_TICKET_CODE, Message: err.Error()})}
		return
	}

	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- event.Chunk{Remaining: false, Data: util.ToJson(core.OnSession{ErrorCode: BAD_REQUEST_CODE, Message: err.Error()})}
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		ch <- event.Chunk{Remaining: false, Data: util.ToJson(core.OnSession{ErrorCode: BAD_REQUEST_CODE, Message: err.Error()})}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ch <- event.Chunk{Remaining: false, Data: util.ToJson(core.OnSession{ErrorCode: BAD_REQUEST_CODE, Message: fmt.Sprintf("http code: %d", resp.StatusCode)})}
		return
	}
	for {
		buff := make([]byte, 1024)
		n, err := resp.Body.Read(buff)
		if n > 0 && err == nil {
			ch <- event.Chunk{Remaining: true, Data: buff[:n]}
			continue
		}
		if n > 0 && err != nil && err == io.EOF {
			ch <- event.Chunk{Remaining: false, Data: buff[:n]}
			break
		}
		if err == io.EOF {
			ch <- event.Chunk{Remaining: false, Data: buff[:0]}
			break
		}
		if err != nil && err != io.EOF {
			ch <- event.Chunk{Remaining: false, Data: buff[:0]}
			core.AppLog.Printf("Resp Error %s\n", err.Error())
			break
		}
	}
}

func (s *AppManager) PostJsonSync(url string, payload any) core.OnSession {
	if s.standalone {
		return core.OnSession{ErrorCode: STANDALONE_APP, Message: STANDALONE_APP_MSG}
	}
	tick, err := s.AppAuth.CreateTicket(1, 1, SUDO_ACCESS_CONTROL)
	if err != nil {
		return core.OnSession{ErrorCode: INVALID_TICKET_CODE, Message: err.Error()}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return core.OnSession{ErrorCode: BAD_REQUEST_CODE, Message: err.Error()}
	}
	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return core.OnSession{ErrorCode: INVALID_JSON_CODE, Message: err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tick)
	resp, err := client.Do(req)
	if err != nil {
		return core.OnSession{ErrorCode: BAD_REQUEST_CODE, Message: err.Error()}
	}
	defer resp.Body.Close()
	var rt core.OnSession
	err = json.NewDecoder(resp.Body).Decode(&rt)
	if err != nil {
		core.AppLog.Printf("Resp Error %s\n", err.Error())
		return core.OnSession{ErrorCode: INVALID_JSON_CODE, Message: err.Error()}
	}
	return rt
}
