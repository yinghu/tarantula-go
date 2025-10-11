package player

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
)

type Simulator struct {
	Player string
	Host   string
}

func (s *Simulator) Play() error {
	err := s.register()
	if err != nil {
		return s.login()
	}
	return nil
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
		hc.Token = session.Token
		hc.Ticket = session.Ticket
		hc.SystemId = session.SystemId
		hc.Home = session.Home
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
		hc.Token = session.Token
		hc.Ticket = session.Ticket
		hc.SystemId = session.SystemId
		hc.Home = session.Home
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
