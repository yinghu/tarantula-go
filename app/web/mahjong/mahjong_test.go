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
	e := event.CreateEvent(cid)
	if e != nil {
		e.OnTopic(topic)
		e.OnListener(s)
		return e, nil
	}
	me := MahjongEvent{}
	me.Callback = s
	return &me, nil
}

func (s *SampleCreator) OnError(e event.Event, err error) {
	fmt.Printf("On event error %v %s\n", e, err.Error())
}

func (s *SampleCreator) OnEvent(e event.Event) {
	fmt.Printf("On event %v\n", e)

}

func TestClient(t *testing.T) {
	hc := util.HttpCaller{Host: "http://192.168.1.11"}
	login := bootstrap.Login{Name: "load1", Hash: "aaa"}
	err := hc.PostJson("presence/register", login, func(resp *http.Response) error {
		session := core.OnSession{}
		err := json.NewDecoder(resp.Body).Decode(&session)
		if err != nil {
			return err
		}
		hc.Token = session.Token
		hc.Ticket = session.Ticket
		hc.SystemId = session.SystemId
		hc.Home = session.Home
		return nil
	})
	if err != nil {
		t.Errorf("login error %s", err.Error())
		return
	}
	fmt.Printf("home %s\n", hc.Home)
	sb := event.SocketPublisher{Remote: fmt.Sprintf("tcp://%s:5050", hc.Home)}
	err = sb.Connect()

	if err != nil {
		t.Errorf("conn error %s", err.Error())
	}

	e := event.JoinEvent{Ticket: hc.Ticket}
	e.OnListener(&SampleCreator{})
	err = sb.Join(&e)
	if err != nil {
		t.Errorf("send error %s", err.Error())
		sb.Close()
		return
	}
	go sb.Subscribe(&SampleCreator{}, &SampleCreator{})

	for range 10 {

		me := MahjongEvent{Cmd: "drop"}
		me.OnTopic("mahjong")
		me.SystemId = hc.SystemId
		me.OnListener(&SampleCreator{})
		sb.Publish(&me, hc.Ticket)
		time.Sleep(10 * time.Millisecond)

	}
	time.Sleep(20 * time.Second)
	sb.Close()
}
