package persistence

import (
	"errors"

	"gameclustering.com/internal/core"
)

const (
	METRICS_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS req_metrics (id BIGSERIAL PRIMARY KEY,path VARCHAR(50) NOT NULL,req_timed BIGINT NOT NULL,req_time TIMESTAMP DEFAULT NOW(),node VARCHAR(10) NOT NULL,req_id INTEGER DEFAULT 0,req_code INTEGER DEFAULT 0)"
	INSERT_METRICS     string = "INSERT INTO req_metrics (path,req_timed,node,req_id,req_code) VALUES($1,$2,$3,$4,$5)"
)

type MetricsDB struct {
	Sql *Postgresql
}

func (s *MetricsDB) Start() error {
	_, err := s.Sql.Exec(METRICS_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil
}

func (s *MetricsDB) WebRequest(m core.ReqMetrics) error {
	inserted, err := s.Sql.Exec(INSERT_METRICS, m.Path, m.ReqTimed, m.Node, m.ReqId, m.ReqCode)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("metrics cannot be saved")
	}
	return nil
}
