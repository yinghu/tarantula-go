package main

import (
	"errors"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/event"
	"github.com/jackc/pgx/v5"
)

func (s *AdminService) createSchema() error {
	_, err := s.Sql.Exec(bootstrap.METRICS_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(bootstrap.ITEM_CONFIGURATION_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(bootstrap.ITEM_HEADER_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(bootstrap.ITEM_APPLICATION_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) SaveLogin(login *event.Login) error {
	inserted, err := s.Sql.Exec("INSERT INTO login (name,hash,access_control) VALUES($1,$2,$3)", login.Name, login.Hash, login.AccessControl)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("login cannot be saved")
	}
	return nil
}

func (s *AdminService) LoadLogin(login *event.Login) error {
	err := s.Sql.Query(func(rows pgx.Rows) error {
		var hash string
		var id int32
		var accessControl int32
		err := rows.Scan(&hash, &id, &accessControl)
		if err != nil {
			return err
		}
		login.Hash = hash
		login.Id = id
		login.AccessControl = accessControl
		return nil
	}, "SELECT hash,id,access_control FROM login WHERE name=$1", login.Name)
	if err != nil {
		return err
	}
	if login.Id == 0 {
		return errors.New("login not existed")
	}
	return nil
}

func (s *AdminService) UpdatePassword(login *event.Login) error {
	updated, err := s.Sql.Exec("UPDATE login SET hash = $1 WHERE name = $2", login.Hash, login.Name)
	if err != nil {
		return err
	}
	if updated == 0 {
		return errors.New("password cannot be saved")
	}
	return nil
}
