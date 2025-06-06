package main

import "gameclustering.com/internal/bootstrap"

func (s *AssetService) createSchema() error {
	_, err := s.Sql.Exec(bootstrap.METRICS_SQL_SCHEMA)
	return err
}
