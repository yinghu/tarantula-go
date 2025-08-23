package main

import (
	"context"
	"fmt"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"github.com/jackc/pgx/v5"
)

const (
	SCHEDULE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_schedule (tournament_id BIGINT PRIMARY KEY,running BOOLEAN NOT NULL)"
	INSTANCE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS tournament_instance (instance_id BIGINT PRIMARY KEY,tournament_id BIGINT NOT NULL,start_time BIGINT NOT NULL,close_time BIGINT NOT NULL,end_time BIGINT NOT NULL,total_entries INTEGER NOT NULL)"
	ENTRY_SQL_SCHEMA    string = "CREATE TABLE IF NOT EXISTS tournament_entry (instance_id BIGINT NOT NULL,system_id BIGINT NOT NULL,score INTEGER NOT NULL,last_updated BIGINT NOT NULL,PRIMARY KEY(instance_id,system_id))"
	JOIN_SQL_SCHEMA     string = "CREATE TABLE IF NOT EXISTS tournament_join (tournament_id BIGINT NOT NULL,system_id BIGINT NOT NULL,instance_id BIGINT NOT NULL,PRIMARY KEY(tournament_id,system_id))"

	INSERT_SCHEDULE string = "INSERT INTO tournament_schedule AS ts (tournament_id,running) VALUES($1,$2) ON CONFLICT (tournament_id) DO UPDATE SET running = true WHERE ts.tournament_id = $3"
	INSERT_INSTANCE string = "INSERT INTO tournament_instance (instance_id,tournament_id,start_time,close_time,end_time,total_entries) VALUES($1,$2,$3,$4,$5,$6)"
	INSERT_JOIN     string = "INSERT INTO tournament_join AS tj (tournament_id,system_id,instance_id) VALUES($1,$2,$3) ON CONFLICT (tournament_id,system_id) DO UPDATE SET instance_id = $4 WHERE tj.tournament_id = $5 AND tj.system_id = $6"
	INSERT_ENTRY    string = "INSERT INTO tournament_entry (instance_id,system_id,score,last_updated) VALUES($1,$2,$3,$4)"

	UPDATE_SCHEDULE string = "UPDATE tournament_schedule AS ts SET running = false WHERE ts.tournament_id = $1"
	UPDATE_SEGMENT  string = "UPDATE tournament_instance AS ti SET total_entries = ti.total_entries + 1 WHERE ti.instance_id = $1 RETURNING total_entries"
	UPDATE_INSTANCE string = "UPDATE tournament_instance AS ti SET total_entries = ti.total_entries + 1 WHERE ti.instance_id = $1 AND ti.total_entries < $2 RETURNING total_entries"
	UPDATE_ENTRY    string = "UPDATE tournament_entry AS te SET score = te.score + $1, last_updated = $2 WHERE te.instance_id = $3 AND te.system_id = $4 RETURNING score"

	SELECT_SCHEDULE string = "SELECT tournament_id FROM tournament_schedule WHERE running = $1"
	SELECT_INSTANCE string = "SELECT instance_id,start_time,close_time,end_time WHERE tournament_id = $1"
	SELECT_JOIN     string = "SELECT instance_id FROM tournament_join WHERE tournament_id = $1 AND system_id = $2"
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
	_, err = s.Sql.Exec(JOIN_SQL_SCHEMA)
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

func (s *TournamentService) updateSegmentSchedule(sc *SegmentSchedule) error {
	r, err := s.Sql.Exec(INSERT_SCHEDULE, sc.TournamentId, true, sc.TournamentId)
	if err != nil {
		return err
	}
	if r == 0 {
		return fmt.Errorf("no row updated")
	}
	for i := range sc.Segments {
		si := sc.Segments[i]
		u, err := s.Sql.Exec(INSERT_INSTANCE, si.InstanceId, sc.TournamentId, sc.StartTime.UnixMilli(), sc.CloseTime.UnixMilli(), sc.EndTime.UnixMilli(), 0)
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
	}, SELECT_SCHEDULE, true)
	if err != nil {
		return ids, err
	}
	return ids, nil
}

func (s *TournamentService) updateInstance(te event.TournamentEvent, limit int32) (int32, error) {
	var total int32
	err := s.Sql.Txn(func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(), UPDATE_INSTANCE, te.InstanceId, limit).Scan(&total)
		if err != nil {
			return nil
		}
		if total == 0 {
			return fmt.Errorf("no row updated")
		}
		return nil
	})
	if err != nil {
		return total, err
	}
	if total == 0 {
		return total, fmt.Errorf("no instance row updated")
	}
	core.AppLog.Printf("Total entries : %d\n", total)
	e, err := s.Sql.Exec(INSERT_ENTRY, te.InstanceId, te.SystemId, 0, time.Now().UnixMilli())
	if err != nil {
		return total, err
	}
	if e == 0 {
		return total, fmt.Errorf("no entry row updated")
	}
	return total, nil
}

func (s *TournamentService) checkJoin(te event.TournamentEvent) event.TournamentEvent {
	s.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&te.InstanceId)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_JOIN, te.TournamentId, te.SystemId)
	return te
}

func (s *TournamentService) updateSegment(te event.TournamentEvent) (int32, error) {
	var total int32
	err := s.Sql.Txn(func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(), UPDATE_SEGMENT, te.InstanceId).Scan(&total)
		if err != nil {
			return nil
		}
		if total == 0 {
			return fmt.Errorf("no row updated")
		}
		r, err := tx.Exec(context.Background(), INSERT_JOIN, te.TournamentId, te.SystemId, te.InstanceId, te.InstanceId, te.TournamentId, te.SystemId)
		if err != nil {
			return err
		}
		if r.RowsAffected() == 0 {
			return fmt.Errorf("no join row updated")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	core.AppLog.Printf("Total entries : %d\n", total)
	e, err := s.Sql.Exec(INSERT_ENTRY, te.InstanceId, te.SystemId, 0, time.Now().UnixMilli())
	if err != nil {
		return 0, err
	}
	if e == 0 {
		return 0, fmt.Errorf("no entry row updated")
	}
	return total, nil
}

func (s *TournamentService) updateEntry(te event.TournamentEvent) (int64, error) {
	var score int64
	err := s.Sql.Txn(func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(), UPDATE_ENTRY, te.Score, te.LastUpdated, te.InstanceId, te.SystemId).Scan(&score)
		if err != nil {
			return nil
		}
		if score == 0 {
			return fmt.Errorf("no row updated")
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return score, nil
}
