package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

const (
	TOPIC_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS topic (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL,app VARCHAR(100) NOT NULL,UNIQUE(name,app))"
	SELECT_TOPICS    string = "SELECT id,name,app FROM topic"
	INSERT_TOPIC     string = "INSERT INTO topic (name,app) VALUES($1,$2) RETURNING id"
	DELETE_TOPIC     string = "DELETE FROM topic WHERE id = $1"
)

func (db *PostofficeService) createSchema() error {
	_, err := db.Sql.Exec(TOPIC_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil
}

func (db *PostofficeService) createTopic(t Topic) (int32, error) {
	var id int32
	err := db.Sql.Txn(func(tx pgx.Tx) error {
		tx.QueryRow(context.Background(), INSERT_TOPIC, t.Name, t.App).Scan(&id)
		return nil
	})
	if err != nil {
		return id, err
	}
	if id == 0 {
		return id, errors.New("no insert")
	}
	return id, nil
}

func (db *PostofficeService) loadTopics() {
	db.Sql.Query(func(row pgx.Rows) error {
		var tp Topic
		err := row.Scan(&tp.Id, &tp.Name, &tp.App)
		if err != nil {
			return err
		}
		db.topics[tp.Id] = tp
		return nil
	}, SELECT_TOPICS)
}
