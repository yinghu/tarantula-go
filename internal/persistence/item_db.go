package persistence

import (
	"context"
	"errors"
	"fmt"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INSERT_CONFIG           string = "INSERT INTO item_configuration (name,type,type_id,category,version) VALUES($1,$2,$3,$4,$5) RETURNING id"
	INSERT_HEADER           string = "INSERT INTO item_header (configuration_id,name,value) VALUES($1,$2,$3)"
	INSERT_APPLICATION      string = "INSERT INTO item_application (configuration_id,name,reference_id) VALUES($1,$2,$3)"
	DELETE_CONFIG_WITH_NAME string = "DELETE FROM item_configuration WHERE name = $1 RETURNING id"
	DELETE_HEADER           string = "DELETE FROM item_header WHERE configuration_id = $1"
	DELETE_APPLICATION      string = "DELETE FROM item_application WHERE configuration_id = $1"
	DELETE_CONFIG_WITH_ID   string = "DELETE FROM item_configuration WHERE id"
	SELECT_CONFIG_WITH_NAME string = "SELECT id,type,type_id,category,version FROM item_configuration WHERE name = $1 LIMIT $2"
)

type ItemDB struct {
	Sql *Postgresql
}

func (db *ItemDB) Save(c item.Configuration) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		var id int32
		err := tx.QueryRow(context.Background(), INSERT_CONFIG, c.Name, c.Type, c.TypeId, c.Category, c.Version).Scan(&id)
		if err != nil {
			return err
		}
		for k, v := range c.Header {
			inserted, err := tx.Exec(context.Background(), INSERT_HEADER, id, k, fmt.Sprintf("%v", v))
			if err != nil {
				return err
			}
			if inserted.RowsAffected() != 1 {
				return errors.New("no data inserted")
			}
		}
		for k, v := range c.Application {
			for i := range v {
				inserted, err := tx.Exec(context.Background(), INSERT_APPLICATION, id, k, v[i])
				if err != nil {
					return err
				}
				if inserted.RowsAffected() != 1 {
					return errors.New("no data inserted")
				}
			}
		}
		return nil
	})
}
func (db *ItemDB) LoadWithName(cname string, limit int) ([]item.Configuration, error) {
	list := make([]item.Configuration, limit)
	ct := 0
	err := db.Sql.Query(func(row pgx.Rows) error {
		conf := item.Configuration{Name: cname}
		err := row.Scan(&conf.Id, &conf.Type, &conf.TypeId, &conf.Category, &conf.Version)
		if err != nil {
			return err
		}
		if conf.Id > 0 {
			fmt.Printf("ID : %d\n", conf.Id)
			list[ct] = conf
			ct++
		}
		return nil
	}, SELECT_CONFIG_WITH_NAME, cname, limit)
	if err != nil {
		return nil, err
	}
	return list[:ct], nil

}

func (db *ItemDB) LoadWithId(cid int32) (item.Configuration, error) {
	return item.Configuration{}, nil
}

func (db *ItemDB) DeleteWithId(cid int32) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), DELETE_CONFIG_WITH_ID, cid)
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), DELETE_HEADER, cid)
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), DELETE_APPLICATION, cid)
		if err != nil {
			return err
		}
		return nil
	})
}

func (db *ItemDB) DeleteWithName(cname string) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		var id int32
		err := tx.QueryRow(context.Background(), DELETE_CONFIG_WITH_NAME, cname).Scan(&id)
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), DELETE_HEADER, id)
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), DELETE_APPLICATION, id)
		if err != nil {
			return err
		}
		return nil
	})
}
