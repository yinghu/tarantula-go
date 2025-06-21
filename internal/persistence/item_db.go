package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	ITEM_ENUM_SQL_SCHEMA       string = "CREATE TABLE IF NOT EXISTS item_enum (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL UNIQUE)"
	ITEM_ENUM_VALUE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_enum_value (enum_id INTEGER,name VARCHAR(100) NOT NULL,value INTEGER NOT NULL,PRIMARY KEY(enum_id,name))"
	ITEM_CATEGORY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_category (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL UNIQUE,scope VARCHAR(30) NOT NULL ,rechargeable BOOL NOT NULL ,description VARCHAR(100) NOT NULL)"
	ITEM_PROPERTY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_category_property (category_id INTEGER NOT NULL,name VARCHAR(100) NOT NULL,type VARCHAR(100) NOT NULL ,reference VARCHAR(100) NOT NULL ,nullable BOOL NOT NULL ,downloadable BOOL NOT NULL, PRIMARY KEY(category_id,name))"

	ITEM_CONFIGURATION_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_configuration (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL,type VARCHAR(50) NOT NULL ,type_id VARCHAR(50) NOT NULL ,category VARCHAR(100) NOT NULL ,version VARCHAR(10) NOT NULL,UNIQUE(name,version))"
	ITEM_HEADER_SQL_SCHEMA        string = "CREATE TABLE IF NOT EXISTS item_header (configuration_id INTEGER NOT NULL,name VARCHAR(100) NOT NULL,value VARCHAR(100) NOT NULL, PRIMARY KEY(configuration_id,name))"
	ITEM_APPLICATION_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_application (configuration_id INTEGER NOT NULL,name VARCHAR(100) NOT NULL,reference_id INTEGER NOT NULL,PRIMARY KEY(configuration_id,name,reference_id))"

	INSERT_ENUM                 string = "INSERT INTO item_enum (name) VALUES ($1) RETURNING id"
	INSERT_ENUM_VALUE           string = "INSERT INTO item_enum_value (enum_id,name,value) VALUES ($1,$2,$3)"
	SELECT_ENUM_WITH_NAME       string = "SELECT id FROM item_enum WHERE name = $1"
	SELECT_ENUM_VALUES_WITH_CID string = "SELECT name,value FROM item_enum_value WHERE enum_id = $1"

	INSERT_CATEGORY            string = "INSERT INTO item_category (name,scope,rechargeable,description) VALUES($1,$2,$3,$4) RETURNING id"
	INSERT_PROPERTY            string = "INSERT INTO item_category_property (category_id,name,type,reference,nullable,downloadable) VALUES($1,$2,$3,$4,$5,$6)"
	SELECT_CATEGORY_WITH_NAME  string = "SELECT id,scope,rechargeable,description FROM item_category WHERE name = $1"
	SELECT_PROPERTIES_WITH_CID string = "SELECT name,type,reference,nullable,downloadable FROM item_category_property WHERE category_id = $1"

	INSERT_CONFIG      string = "INSERT INTO item_configuration (name,type,type_id,category,version) VALUES($1,$2,$3,$4,$5) RETURNING id"
	INSERT_HEADER      string = "INSERT INTO item_header (configuration_id,name,value) VALUES($1,$2,$3)"
	INSERT_APPLICATION string = "INSERT INTO item_application (configuration_id,name,reference_id) VALUES($1,$2,$3)"

	DELETE_CONFIG_WITH_NAME string = "DELETE FROM item_configuration WHERE name = $1 RETURNING id"
	DELETE_HEADER           string = "DELETE FROM item_header WHERE configuration_id = $1"
	DELETE_APPLICATION      string = "DELETE FROM item_application WHERE configuration_id = $1"
	DELETE_CONFIG_WITH_ID   string = "DELETE FROM item_configuration WHERE id"

	SELECT_CONFIG_WITH_NAME           string = "SELECT id,type,type_id,category,version FROM item_configuration WHERE name = $1 LIMIT $2"
	SELECT_CONFIG_WITH_ID             string = "SELECT name,type,type_id,category,version FROM item_configuration WHERE id = $1"
	SELECT_CONFIG_HEADER_WIHT_ID      string = "SELECT name,value FROM item_header WHERE configuration_id = $1"
	SELECT_CONFIG_APPLICATION_WITH_ID string = "SELECT name,reference_id FROM item_application WHERE configuration_id = $1"
)

type ItemDB struct {
	Sql *Postgresql
}

func (db *ItemDB) Start() error {
	_, err := db.Sql.Exec(ITEM_ENUM_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = db.Sql.Exec(ITEM_ENUM_VALUE_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = db.Sql.Exec(ITEM_CATEGORY_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = db.Sql.Exec(ITEM_PROPERTY_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = db.Sql.Exec(ITEM_CONFIGURATION_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = db.Sql.Exec(ITEM_HEADER_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = db.Sql.Exec(ITEM_APPLICATION_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil
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
		return cat, errors.New("category not existed")
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

func (db *ItemDB) SaveEnum(c item.Enum) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		var id int32
		err := tx.QueryRow(context.Background(), INSERT_ENUM, c.Name).Scan(&id)
		if err != nil {
			return err
		}
		valid := false
		for i := range c.Values {
			p := c.Values[i]
			_, err = tx.Exec(context.Background(), INSERT_ENUM_VALUE, id, p.Name, p.Value)
			if err != nil {
				valid = false
				return err
			}
			valid = true
		}
		if !valid {
			return errors.New("at least 1 enum value required")
		}
		return nil
	})
}

func (db *ItemDB) LoadEnum(cname string) (item.Enum, error) {
	cat := item.Enum{Name: cname}
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&cat.Id)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_ENUM_WITH_NAME, cname)
	if err != nil {
		return cat, err
	}
	if cat.Id == 0 {
		return cat, errors.New("enum not existed")
	}
	cat.Values = make([]item.EnumValue, 0)
	err = db.Sql.Query(func(row pgx.Rows) error {
		var prop item.EnumValue
		err := row.Scan(&prop.Name, &prop.Value)
		if err != nil {
			return err
		}
		cat.Values = append(cat.Values, prop)
		return nil
	}, SELECT_ENUM_VALUES_WITH_CID, cat.Id)
	if err != nil {
		return cat, errors.New("no enum value existed")
	}
	return cat, nil
}

func (db *ItemDB) ValidateEnum(c item.Enum) error {
	if c.Name == "" {
		return errors.New("name none empty string required")
	}
	if len(c.Values) == 0 {
		return errors.New("at least 1 value required")
	}
	for i := range c.Values {
		v := c.Values[i]
		if v.Name == "" {
			return errors.New("prop name none empty string required")
		}
	}
	return nil
}

func (db *ItemDB) ValidateCategory(c item.Category) error {
	if c.Scope == "" {
		return errors.New("scope none empty string required")
	}
	if c.Name == "" {
		return errors.New("name none empty string required")
	}
	if c.Description == "" {
		return errors.New("description none empty string required")
	}
	if len(c.Properties) == 0 {
		return errors.New("at least 1 property required")
	}
	for i := range c.Properties {
		prop := c.Properties[i]
		if prop.Name == "" {
			return errors.New("prop name none empty string required")
		}
		if prop.Type == "" {
			return errors.New("prop type none empty string required")
		}
		if prop.Reference == "" {
			return errors.New("prop reference none empty string required")
		}
		if prop.Type == "enum" {
			_, err := db.LoadEnum(prop.Reference)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "category" {
			parts := strings.Split(prop.Reference, ":")
			if len(parts) != 2 {
				return errors.New("wrong category reference format")
			}
			_, err := db.LoadCategory(parts[1])
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "set" || prop.Type == "list" {
			parts := strings.Split(prop.Reference, ":")
			if len(parts) != 2 {
				return errors.New("wrong category reference format")
			}
			_, err := db.LoadCategory(parts[1])
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
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
		if prop.Type == "category" || prop.Type == "set" || prop.Type == "list" {
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
		if prop.Type == "number" && prop.Reference == "int" {
			err = asInt(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "number" && prop.Reference == "long" {
			err = asLong(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "number" && prop.Reference == "float" {
			err = asFloat(v)
			if err != nil {
				return err
			}
			continue
		}
		if prop.Type == "number" && prop.Reference == "double" {
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
			_, err = time.Parse(prop.Reference, fmt.Sprintf("%v", v))
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
