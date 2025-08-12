package main

const (
	SCHEDULE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_schedule (tournament_id BIGINT PRIMARY KEY,running BOOLEAN NOT NULL)"
	INSTANCE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_instance (instance_id BIGINT PRIMARY KEY,tournament_id BIGINT NOT NULL,start_time TIMESTAMP NOT NULL,close_time TIMESTAMP NOT NULL,end_time TIMESTAMP NOT NULL)"
	ENTRY_SQL_SCHEMA    string = "CREATE TABLE IF NOT EXISTS tournament_entry (instance_id BIGINT NOT NULL,system_id BIGINT NOT NULL,score INTEGER NOT NULL,last_updated TIMESTAMP NOT NULL,PRIMARY KEY(instance_id,system_id))"
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
