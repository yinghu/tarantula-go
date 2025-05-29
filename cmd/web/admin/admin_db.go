package main

import (
	"errors"

	"gameclustering.com/internal/event"
	"github.com/jackc/pgx/v5"
)

func (s *AdminService) SaveLogin(login *event.Login) error {
	inserted, err := s.sql.Exec("INSERT INTO login (name,hash) VALUES($1,$2)", login.Name, login.Hash)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("login cannot be saved")
	}
	return nil
}

func (s *AdminService) LoadLogin(login *event.Login) error {
	err := s.sql.Query(func(rows pgx.Rows) error {
		var hash string
		var systemId int64
		err := rows.Scan(&hash, &systemId)
		if err != nil {
			return err
		}
		login.Hash = hash
		login.SystemId = systemId
		return nil
	}, "SELECT hash,id FROM login WHERE name=$1", login.Name)
	if err != nil {
		return err
	}
	if login.SystemId == 0 {
		return errors.New("login not existed")
	}
	return nil
}


