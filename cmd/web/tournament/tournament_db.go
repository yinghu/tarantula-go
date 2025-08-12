package main

import (
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
	SCHEDULE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_schedule (tournament_id BIGINT PRIMARY KEY,running BOOLEAN NOT NULL)"
	INSTANCE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_instance (instance_id BIGINT PRIMARY KEY,tournament_id BIGINT NOT NULL,start_time BIGINT NOT NULL,close_time BIGINT NOT NULL,end_time BIGINT NOT NULL)"
	ENTRY_SQL_SCHEMA    string = "CREATE TABLE IF NOT EXISTS tournament_entry (instance_id BIGINT NOT NULL,system_id BIGINT NOT NULL,score INTEGER NOT NULL,last_updated BIGINT NOT NULL,PRIMARY KEY(instance_id,system_id))"

	INSERT_SCHEDULE string = "INSERT INTO tournament_schedule AS ts (tournament_id,running) VALUES($1,$2) ON CONFLICT (tournament_id) DO UPDATE SET running = true WHERE ts.tournament_id = $3"
	INSERT_INSTANCE string = "INSERT INTO tournament_instance (instance_id,tournament_id,start_time,close_time,end_time) VALUES($1,$2,$3,$4,$5)"

	UPDATE_SCHEDULE string = "UPDATE tournament_schedule AS ts SET running = false WHERE ts.tournament_id = $1"

	SELECT_SCHEDULE string = "SELECT tournament_id FROM tournament_schedule WHERE running = $1"
	SELECT_INSTANCE string = "SELECT instance_id,start_time,close_time,end_time WHERE tournament_id = $1"
)

func (s *TournamentService) createSchema() error {
	_, err := s.Sql.Exec(SCHEDULE_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(INSTANCE_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(ENTRY_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil
}

func (s *TournamentService) updateInstanceSchedule(sc InstanceSchedule) error {
	r, err := s.Sql.Exec(INSERT_SCHEDULE, sc.TournamentId, true, sc.TournamentId)
	if err != nil {
		return err
	}
	if r == 0 {
		return fmt.Errorf("no row updated")
	}
	return nil
}

func (s *TournamentService) updateSegmentSchedule(sc SegementSchedule) error {
	r, err := s.Sql.Exec(INSERT_SCHEDULE, sc.TournamentId, true, sc.TournamentId)
	if err != nil {
		return err
	}
	if r == 0 {
		return fmt.Errorf("no row updated")
	}
	for i := range sc.Segments {
		si := sc.Segments[i]
		u, err := s.Sql.Exec(INSERT_INSTANCE, si.InstanceId, sc.TournamentId, sc.StartTime.UnixMilli(), sc.CloseTime.UnixMilli(), sc.EndTime.UnixMilli())
		if err != nil {
			return err
		}
		if u == 0 {
			return fmt.Errorf("no row updated")
		}
	}
	return nil
}
func (s *TournamentService) updateSchedule(id int64) error {
	r, err := s.Sql.Exec(UPDATE_SCHEDULE, id)
	if err != nil {
		return err
	}
	if r == 0 {
		return fmt.Errorf("no row updated")
	}
	return nil
}
func (s *TournamentService) loadSchedule() ([]int64, error) {
	ids := make([]int64, 0)
	err := s.Sql.Query(func(row pgx.Rows) error {
		var id int64
		err := row.Scan(&id)
		if err != nil {
			return err
		}
		if id > 0 {
			ids = append(ids, id)
		}
		return nil
	}, SELECT_SCHEDULE,true)
	if err != nil {
		return ids, err
	}
	return ids, nil
}
