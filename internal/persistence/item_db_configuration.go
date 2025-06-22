package persistence

import (
	"context"
	"errors"
	"fmt"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

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
			list[ct] = conf
			ct++
		}
		return nil
	}, SELECT_CONFIG_WITH_NAME, cname, limit)
	if err != nil {
		return nil, err
	}
	for i := range list[:ct] {
		db.loadHeader(&list[i])
		db.loadApplication(&list[i])
	}
	return list[:ct], nil

}

func (db *ItemDB) LoadWithId(cid int32) (item.Configuration, error) {
	conf := item.Configuration{Id: cid}
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&conf.Name, &conf.Type, &conf.TypeId, &conf.Category, &conf.Version)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_CONFIG_WITH_ID, cid)
	if err != nil {
		return conf, err
	}
	if conf.Name == "" {
		return conf, errors.New("obj not existed")
	}
	return conf, nil
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

func (db *ItemDB) loadHeader(c *item.Configuration) error {
	c.Header = make(map[string]any)
	return db.Sql.Query(func(row pgx.Rows) error {
		var k string
		var v any
		err := row.Scan(&k, &v)
		if err != nil {
			return err
		}
		c.Header[k] = v
		return nil
	}, SELECT_CONFIG_HEADER_WIHT_ID, c.Id)
}

func (db *ItemDB) loadApplication(c *item.Configuration) error {
	c.Application = make(map[string][]int32)
	return db.Sql.Query(func(row pgx.Rows) error {
		var k string
		var v int32
		err := row.Scan(&k, &v)
		if err != nil {
			return err
		}
		c.Application[k] = append(c.Application[k], v)
		return nil
	}, SELECT_CONFIG_APPLICATION_WITH_ID, c.Id)
}
