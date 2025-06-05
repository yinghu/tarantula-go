package persistence

import (
	"errors"

	"gameclustering.com/internal/metrics"
)

type MetricsDB struct {
	Sql *Postgresql
}

func (s *MetricsDB) WebRequest(m metrics.ReqMetrics) error {
	inserted, err := s.Sql.Exec("INSERT INTO req_metrics (path,req_timed,node,req_id) VALUES($1,$2,$3,$4)", m.Path, m.ReqTimed, m.Node, m.ReqId)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("metrics cannot be saved")
	}
	return nil
}
