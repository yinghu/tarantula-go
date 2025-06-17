package persistence

import (
	"context"
	"errors"
	"fmt"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INSERT_CONFIG string = "INSERT INTO item_configuration (name,type,type_id,category,version) VALUES($1,$2,$3,$4,$5) RETURNING id"
	INSERT_HEADER string = "INSERT INTO item_header (configuration_id,name,value) VALUES($1,$2,$3)"
)

type ItemDB struct {
	Sql *Postgresql
}

func (db *ItemDB) Save(c item.Configuration) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		var id int32
		err := tx.QueryRow(context.Background(), INSERT_CONFIG, c.Name, c.Type, c.TypeId, c.Category, c.Version).Scan(&id)
		if err != nil {
			fmt.Printf("Error %s\n", err.Error())
			return err
		}
		fmt.Printf("id : %d\n", id)
		for k, v := range c.Header {

			inserted, err := tx.Exec(context.Background(), INSERT_HEADER, id, k, fmt.Sprintf("%v", v))
			if err != nil {
				fmt.Printf("Error %s\n", err.Error())
				return err
			}
			if inserted.RowsAffected() != 1 {
				return errors.New("no data inserted")
			}
		}
		return nil

	})
}
func (db *ItemDB) LoadWithName(cname string) (item.Configuration, error) {
	return item.Configuration{}, nil
}

func (db *ItemDB) LoadWithId(cid int32) (item.Configuration, error) {
	return item.Configuration{}, nil
}
