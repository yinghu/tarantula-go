package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

type SampleCreator struct {
}

func (s *SampleCreator) Create(cid int, topic string) (event.Event, error) {
	me := MahjongEvent{}
	me.Callback = &MahjongEventListener{}
	return &me, nil
}

func TestClient(t *testing.T) {
	hc := util.HttpCaller{Host: "http://192.168.1.11"}
	login := bootstrap.Login{Name: "player1", Hash: "aaa"}
	err := hc.PostJson("presence/login", login, func(resp *http.Response) error {
		session := core.OnSession{}
		err := json.NewDecoder(resp.Body).Decode(&session)
		if err != nil {
			return err
		}
		hc.Token = session.Token
		return nil
	})
	if err != nil {
		t.Errorf("login error %s", err.Error())
		return
	}
	fmt.Printf("token %s\n",hc.Token)
	sb := event.SocketPublisher{Remote: "tcp://192.168.1.11:5050"}
	err = sb.Connect()
	if err != nil {
		t.Errorf("conn error %s", err.Error())
	}

	e := event.JoinEvent{Token: hc.Token}
	e.OnListener(&MahjongEventListener{})
	err = sb.Join(&e)
	if err != nil {
		t.Errorf("send error %s", err.Error())
		sb.Close()
		return
	}
	go sb.Subscribe(&SampleCreator{}, &MahjongEventListener{})
	for range 3 {
		me := MahjongEvent{Cmd: "drop"}
		me.OnListener(&MahjongEventListener{})
		sb.Publish(&me,hc.Token)
	}
	time.Sleep(5 * time.Second)
	sb.Close()
}
