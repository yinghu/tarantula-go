package player

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/util"
	"gameclustering.com/internal/event"
)

type Simulator struct {
	Player   string
	Host     string
	Token    string
	Ticket   string
	SystemId int64
	Home     string
}

func (s *Simulator) Play() error {
	err := s.register()
	if err != nil {
		err = s.login()
		if err != nil {
			return err
		}
	}
	err = s.inventory()
	if err != nil {
		return err
	}
	done := make(chan bool)
	go s.tcp(done)
	b := <-done
	if b {
		return nil
	}
	return fmt.Errorf("tcp error")
}

func (s *Simulator) register() error {
	hc := util.HttpCaller{Host: s.Host}
	login := bootstrap.Login{Name: s.Player, Hash: "password"}
	err := hc.PostJson("presence/register", login, func(resp *http.Response) error {
		session := core.OnSession{}
		err := json.NewDecoder(resp.Body).Decode(&session)
		if err != nil {
			return err
		}
		if !session.Successful {
			return fmt.Errorf("error : %s", session.Message)
		}
		s.Token = session.Token
		s.Ticket = session.Ticket
		s.SystemId = session.SystemId
		s.Home = session.Home
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Simulator) login() error {
	hc := util.HttpCaller{Host: s.Host}
	login := bootstrap.Login{Name: s.Player, Hash: "password"}
	err := hc.PostJson("presence/login", login, func(resp *http.Response) error {
		session := core.OnSession{}
		err := json.NewDecoder(resp.Body).Decode(&session)
		if err != nil {
			return err
		}
		if !session.Successful {
			return fmt.Errorf("error : %s", session.Message)
		}
		s.Token = session.Token
		s.Ticket = session.Ticket
		s.SystemId = session.SystemId
		s.Home = session.Home
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Simulator) inventory() error {
	hc := util.HttpCaller{Host: s.Host, Token: s.Token, SystemId: s.SystemId}
	req := item.OnInventory{SystemId: hc.SystemId, TypeId: "gold"}
	err := hc.PostJson("inventory/load", req, func(resp *http.Response) error {
		inv := persistence.InventoryResp{}
		err := json.NewDecoder(resp.Body).Decode(&inv)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Simulator) tcp(ch chan bool) {
	sb := event.TcpPublisher{Remote: fmt.Sprintf("tcp://%s:5050", s.Home)}
	err := sb.Connect()

	if err != nil {
		ch <- false
	}

	e := event.JoinEvent{Ticket: s.Ticket}
	e.OnListener(&SampleCreator{})
	err = sb.Join(&e)
	if err != nil {
		ch <- false
		sb.Close()
		return
	}
	go sb.Subscribe(&SampleCreator{}, &SampleCreator{})

	for range 10 {
		me := MahjongEvent{Cmd: 0}
		me.OnTopic("mahjong")
		me.SystemId = s.SystemId
		me.OnListener(&SampleCreator{})
		sb.Publish(&me, s.Ticket)
		time.Sleep(1000 * time.Millisecond)
	}
	time.Sleep(1 * time.Second)
	sb.Close()
	ch <- true
}
