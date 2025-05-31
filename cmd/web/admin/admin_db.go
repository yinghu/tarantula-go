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
		var accessControl int32
		err := rows.Scan(&hash, &systemId, &accessControl)
		if err != nil {
			return err
		}
		login.Hash = hash
		login.SystemId = systemId
		login.AccessControl = accessControl
		return nil
	}, "SELECT hash,id,access_control FROM login WHERE name=$1", login.Name)
	if err != nil {
		return err
	}
	if login.SystemId == 0 {
		return errors.New("login not existed")
	}
	return nil
}
