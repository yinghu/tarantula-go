package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	DATE_TIME_FORMAT string = "2006-01-02T15:04"

	INSERT_CONFIG      string = "INSERT INTO item_configuration (id,name,type,type_id,category,version) VALUES($1,$2,$3,$4,$5,$6)"
	INSERT_HEADER      string = "INSERT INTO item_header (configuration_id,name,value) VALUES($1,$2,$3)"
	INSERT_APPLICATION string = "INSERT INTO item_application (configuration_id,name,reference_id) VALUES($1,$2,$3)"

	SELECT_CONFIG_WITH_NAME           string = "SELECT id,name,type,type_id,version FROM item_configuration WHERE category = $1 LIMIT $2"
	SELECT_CONFIG_WITH_ID             string = "SELECT name,type,type_id,category,version FROM item_configuration WHERE id = $1"
	SELECT_CONFIG_HEADER_WIHT_ID      string = "SELECT name,value FROM item_header WHERE configuration_id = $1"
	SELECT_CONFIG_APPLICATION_WITH_ID string = "SELECT name,reference_id FROM item_application WHERE configuration_id = $1"
)

func (db *ItemDB) Save(c item.Configuration) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		r, err := tx.Exec(context.Background(), INSERT_CONFIG, c.Id, c.Name, c.Type, c.TypeId, c.Category, c.Version)
		if err != nil {
			return err
		}
		if r.RowsAffected() == 0 {
			return errors.New("no insert")
		}
		for k, v := range c.Header {
			inserted, err := tx.Exec(context.Background(), INSERT_HEADER, c.Id, k, fmt.Sprintf("%v", v))
			if err != nil {
				return err
			}
			if inserted.RowsAffected() != 1 {
				return errors.New("no data inserted")
			}
		}
		for k, v := range c.Application {
			for i := range v {
				inserted, err := tx.Exec(context.Background(), INSERT_APPLICATION, c.Id, k, v[i])
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
		conf := item.Configuration{Category: cname}
		err := row.Scan(&conf.Id, &conf.Name, &conf.Type, &conf.TypeId, &conf.Version)
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

func (db *ItemDB) LoadWithId(cid int64) (item.Configuration, error) {
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
	db.loadHeader(&conf)
	db.loadApplication(&conf)
	return conf, nil
}

func (db *ItemDB) DeleteWithId(cid int64) error {
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
	c.Application = make(map[string][]int64)
	return db.Sql.Query(func(row pgx.Rows) error {
		var k string
		var v int64
		err := row.Scan(&k, &v)
		if err != nil {
			return err
		}
		c.Application[k] = append(c.Application[k], v)
		return nil
	}, SELECT_CONFIG_APPLICATION_WITH_ID, c.Id)
}

func (db *ItemDB) Validate(c item.Configuration) error {
	if c.Name == "" {
		return errors.New("name none empty string required")
	}
	if c.TypeId == "" {
		return errors.New("typeId none empty string required")
	}
	if c.Type == "" {
		return errors.New("type none empty string required")
	}
	if c.Category == "" {
		return errors.New("category none empty string required")
	}
	if c.Version == "" {
		return errors.New("version none empty string required")
	}
	cat, err := db.LoadCategory(c.Category)
	if err != nil {
		return err
	}
	valid := len(c.Header)
	for i := range cat.Properties {
		prop := cat.Properties[i]
		if prop.Type == "category" || prop.Type == "set" || prop.Type == "list" || prop.Type == "scope" {
			for _, v := range c.Application {
				for i := range v {
					_, err := db.LoadWithId(v[i])
					if err != nil {
						return err
					}
				}
			}
			continue
		}
		v, existed := c.Header[prop.Name]
		if !existed && !prop.Nullable {
			return errors.New("value not existed : " + prop.Type)
		}
		valid--
		if prop.Type == "string" {
			err = asString(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "int" {
			err = asInt(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "long" {
			err = asLong(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "float" {
			err = asFloat(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "double" {
			err = asDouble(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "boolean" {
			err = asBool(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "dateTime" {
			err = asString(v)
			if err != nil {
				return err
			}
			_, err = time.Parse(DATE_TIME_FORMAT, fmt.Sprintf("%v", v))
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "enum" {
			em, err := db.LoadEnum(prop.Reference)
			if err != nil {
				return err
			}
			e, err := toInt32(v)
			if err != nil {
				return err
			}
			matched := false
			for i := range em.Values {
				matched = em.Values[i].Value == e
				if matched {
					break
				}
			}
			if !matched {
				return errors.New("enum value not matched")
			}
			continue
		}
	}
	if valid == 0 {
		return nil
	}
	return errors.New("invalid data")
}

func asString(v any) error {
	_, ok := v.(string)
	if ok {
		return nil
	}
	return errors.New("wrong string format")
}

func asDouble(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong double format")
}

func asFloat(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong float format")
}

func asInt(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong int format")
}

func toInt32(v any) (int32, error) {
	x, ok := v.(float64)
	if ok {
		return int32(x), nil
	}
	return 0, errors.New("wrong int format")
}

func asLong(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong long format")
}

func asBool(v any) error {
	_, ok := v.(bool)
	if ok {
		return nil
	}
	return errors.New("wrong bool format")
}
