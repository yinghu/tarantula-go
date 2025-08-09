package main

const (
	TOURNAMENT_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament (id SERIAL PRIMARY KEY,tournament_id BIGINT NOT NULL,instance_id BIGINT NOT NULL,system_id BIGINT NOT NULL,score INTEGER NOT NULL,last_updated TIMESTAMP NOT NULL,UNIQUE(tournament_id,instance_id,system_id))"
)

func (s *TournamentService) createSchema() error {
	_, err := s.Sql.Exec(TOURNAMENT_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil
}
