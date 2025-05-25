package main

import (
	"errors"

	"gameclustering.com/internal/event"
	"gameclustering.com/internal/metrics"
	"github.com/jackc/pgx/v5"
)

func (s *PresenceService) SaveLogin(login *event.Login) error {
	inserted, err := s.Sql.Exec("INSERT INTO login (name,hash,system_id,reference_id) VALUES($1,$2,$3,$4)", login.Name, login.Hash, login.SystemId, login.ReferenceId)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("login cannot be saved")
	}
	return nil
}

func (s *PresenceService) LoadLogin(login *event.Login) error {
	err := s.Sql.Query(func(rows pgx.Rows) error {
		var hash string
		var systemId int64
		err := rows.Scan(&hash, &systemId)
		if err != nil {
			return err
		}
		login.Hash = hash
		login.SystemId = systemId
		return nil
	}, "SELECT hash,system_id FROM login WHERE name=$1", login.Name)
	if err != nil {
		return err
	}
	if login.SystemId == 0 {
		return errors.New("login not existed")
	}
	return nil
}

func (s *PresenceService) SaveMetrics(metrics *metrics.ReqMetrics) error {
	inserted, err := s.Sql.Exec("INSERT INTO metrics (path,req_timed,node) VALUES($1,$2,$3)", metrics.Path, metrics.ReqTimed, metrics.Node)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("metrics cannot be saved")
	}
	return nil
}
