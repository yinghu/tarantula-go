package main

import "fmt"

const (
	SCHEDULE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_schedule (tournament_id BIGINT PRIMARY KEY,running BOOLEAN NOT NULL)"
	INSTANCE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_instance (instance_id BIGINT PRIMARY KEY,tournament_id BIGINT NOT NULL,start_time TIMESTAMP NOT NULL,close_time TIMESTAMP NOT NULL,end_time TIMESTAMP NOT NULL)"
	ENTRY_SQL_SCHEMA    string = "CREATE TABLE IF NOT EXISTS tournament_entry (instance_id BIGINT NOT NULL,system_id BIGINT NOT NULL,score INTEGER NOT NULL,last_updated TIMESTAMP NOT NULL,PRIMARY KEY(instance_id,system_id))"

	INSERT_SCHEDULE string = "INSERT INTO tournament_schedule AS ts (tournament_id,running) VALUES($1,$2) ON CONFLICT DO UPDATE SET ts.running = true WHERE ts.tournament_id = $3"
	INSERT_INSTANCE string = "INSERT INTO tournament_instance (instance_id,tournament_id,start_time,close_time,end_time) VALUES($1,$2,$3,$4,$5)"

	UPDATE_SCHEDULE string = "UPDATE tournament_schedule AS ts SET ts.running = false WHERE ts.tournament_id = $1"

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
