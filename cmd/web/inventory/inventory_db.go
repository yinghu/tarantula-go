package main

import (
	//"errors"

	"gameclustering.com/internal/bootstrap"
	//"github.com/jackc/pgx/v5"
)

func (s *InventoryService) createSchema() error {
	_, err := s.Sql.Exec(bootstrap.METRICS_SQL_SCHEMA)
	if err != nil {
		return err
	}

	return nil
}
