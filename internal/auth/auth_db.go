package auth

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Service) SaveLogin(login *Login) error {
	inserted, err := s.Sql.Exec("INSERT INTO login (name,hash,system_id,reference_id) VALUES($1,$2,$3,$4)", login.Name, login.Hash, login.SystemId, login.ReferenceId)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("login cannot be saved")
	}
	return nil
}

func (s *Service) LoadLogin(login *Login) error {
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
