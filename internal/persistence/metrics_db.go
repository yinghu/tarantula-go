package persistence

import (
	"errors"
	"gameclustering.com/internal/metrics"
)

type MetricsDB struct {
	Sql *Postgresql
}

func (s *MetricsDB) WebRequest(m metrics.ReqMetrics) error {
	inserted, err := s.Sql.Exec("INSERT INTO req_metrics (path,req_timed,node) VALUES($1,$2,$3)", m.Path, m.ReqTimed, m.Node)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("metrics cannot be saved")
	}
	return nil
}

