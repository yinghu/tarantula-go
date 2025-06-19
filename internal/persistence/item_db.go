package persistence

import (
	"context"
	"errors"
	"fmt"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INSERT_CATEGORY            string = "INSERT INTO item_category (name,scope,rechargeable,description) VALUES($1,$2,$3,$4) RETURNING id"
	INSERT_PROPERTY            string = "INSERT INTO item_property (category_id,name,type,reference,nullable,downloadable) VALUES($1,$2,$3,$4,$5,$6)"
	SELECT_CATEGORY_WITH_NAME  string = "SELECT id,scope,rechargeable,description FROM item_category WHERE name = $1"
	SELECT_PROPERTIES_WITH_CID string = "SELECT name,type,reference,nullable,downloadable FROM item_property WHERE category_id = $1"

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

func (db *ItemDB) SaveCategory(c item.Category) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		var id int32
		err := tx.QueryRow(context.Background(), INSERT_CATEGORY, c.Name, c.Scope, c.Rechargeable, c.Description).Scan(&id)
		if err != nil {
			return err
		}
		valid := false
		for i := range c.Properties {
			p := c.Properties[i]
			_, err = tx.Exec(context.Background(), INSERT_PROPERTY, id, p.Name, p.Type, p.Reference, p.Nullable, p.Downloadable)
			if err != nil {
				valid = false
				return err
			}
			valid = true
		}
		if !valid {
			return errors.New("at least 1 property required")
		}
		return nil
	})
}

func (db *ItemDB) LoadCategory(cname string) (item.Category, error) {
	cat := item.Category{Name: cname}
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&cat.Id, &cat.Scope, &cat.Rechargeable, &cat.Description)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_CATEGORY_WITH_NAME, cname)
	if err != nil {
		return cat, err
	}
	if cat.Id == 0 {
		return cat, errors.New("not existed")
	}
	cat.Properties = make([]item.Property, 0)
	err = db.Sql.Query(func(row pgx.Rows) error {
		var prop item.Property
		err := row.Scan(&prop.Name, &prop.Type, &prop.Reference, &prop.Nullable, &prop.Downloadable)
		if err != nil {
			return err
		}
		cat.Properties = append(cat.Properties, prop)
		return nil
	}, SELECT_PROPERTIES_WITH_CID, cat.Id)
	if err != nil {
		return cat, errors.New("no property existed")
	}
	return cat, nil
}

func (db *ItemDB) Validate(c item.Configuration) error {
	cat, err := db.LoadCategory(c.Category)
	if err != nil {
		return err
	}
	for i := range cat.Properties {
		prop := cat.Properties[i]
		_, existed := c.Header[prop.Name]
		if !existed {
			return errors.New(prop.Name + " not existed")
		}

	}
	return nil
}
