package main

import (
	"errors"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/metrics"
	"github.com/jackc/pgx/v5"
)

const (
	CREATE_LOGIN_SCHEMA    string = "CREATE TABLE IF NOT EXISTS login (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL UNIQUE,hash VARCHAR(255) NOT NULL,system_id BIGINT NOT NULL UNIQUE,reference_id INTEGER DEFAULT 0,access_control INTEGER DEFAULT 1)"
	INSERT_LOGIN           string = "INSERT INTO login (name,hash,system_id,reference_id) VALUES($1,$2,$3,$4)"
	SELECT_LOGIN_WITH_NAME        = "SELECT hash,system_id,access_control,id FROM login WHERE name=$1"
	UPDATE_HASH                   = "UPDATE login SET hash = $1 WHERE name = $2"
)

func (s *PresenceService) createSchema() error {
	_, err := s.Sql.Exec(CREATE_LOGIN_SCHEMA)
	if err != nil {
		return err
	}
	return nil
}

func (s *PresenceService) SaveLogin(login *bootstrap.Login) error {
	inserted, err := s.Sql.Exec(INSERT_LOGIN, login.Name, login.Hash, login.SystemId, login.ReferenceId)
	if err != nil {
		metrics.APP_ERROR_METRICS.WithLabelValues("login", err.Error()).Inc()
		return err
	}
	if inserted == 0 {
		metrics.APP_ERROR_METRICS.WithLabelValues("login", "cannot be saved").Inc()
		return errors.New("login cannot be saved")
	}
	return nil
}

func (s *PresenceService) LoadLogin(login *bootstrap.Login) error {
	err := s.Sql.Query(func(rows pgx.Rows) error {
		var hash string
		var systemId int64
		var accessControl int32
		var id int32
		err := rows.Scan(&hash, &systemId, &accessControl, &id)
		if err != nil {
			return err
		}
		login.Hash = hash
		login.SystemId = systemId
		login.AccessControl = accessControl
		login.Id = id
		return nil
	}, SELECT_LOGIN_WITH_NAME, login.Name)
	if err != nil {
		return err
	}
	if login.SystemId == 0 {
		return errors.New("login not existed")
	}
	return nil
}

func (s *PresenceService) UpdatePassword(login *bootstrap.Login) error {
	updated, err := s.Sql.Exec(UPDATE_HASH, login.Hash, login.Name)
	if err != nil {
		return err
	}
	if updated == 0 {
		return errors.New("password cannot be saved")
	}
	return nil
}
