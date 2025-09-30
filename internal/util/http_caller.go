package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Callback func(resp *http.Response) error

type HttpCaller struct {
	Host     string
	SystemId int64
	Token    string
	Ticket   string
	Home     string
}

func (h *HttpCaller) PostJson(path string, payload any, cb Callback) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", h.Host, path), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if h.Token != "" {
		req.Header.Set("Authorization", "Bearer "+h.Token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return cb(resp)
}
func (h *HttpCaller) GetJson(path string, cb Callback) error {

	tr := &http.Transport{
		DisableKeepAlives:  true,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", h.Host, path), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if h.Token != "" {
		req.Header.Set("Authorization", "Bearer "+h.Token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return cb(resp)
}
