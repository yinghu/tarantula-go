package persistence

import (
	"context"
	"errors"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INSERT_ENUM                 string = "INSERT INTO item_enum (id,name) VALUES ($1,$2)"
	INSERT_ENUM_VALUE           string = "INSERT INTO item_enum_value (enum_id,name,value) VALUES ($1,$2,$3)"
	SELECT_ENUM_WITH_NAME       string = "SELECT id FROM item_enum WHERE name = $1"
	SELECT_ALL_ENUM             string = "SELECT id,name FROM item_enum"
	SELECT_ENUM_VALUES_WITH_CID string = "SELECT name,value FROM item_enum_value WHERE enum_id = $1"
)

func (db *ItemDB) SaveEnum(c item.Enum) error {
	err := db.validateEnum(c)
	if err != nil {
		return err
	}
	return db.Sql.Txn(func(tx pgx.Tx) error {
		r, err := tx.Exec(context.Background(), INSERT_ENUM, c.Id, c.Name)
		if err != nil {
			return err
		}
		if r.RowsAffected() == 0 {
			return errors.New("no insert")
		}
		valid := false
		for i := range c.Values {
			p := c.Values[i]
			_, err = tx.Exec(context.Background(), INSERT_ENUM_VALUE, c.Id, p.Name, p.Value)
			if err != nil {
				valid = false
				return err
			}
			valid = true
		}
		if !valid {
			return errors.New("at least 1 enum value required")
		}
		return db.Gis.SaveEnum(c)
	})
}

func (db *ItemDB) validateEnum(c item.Enum) error {
	if c.Id <= 0 {
		return errors.New("none negative id required")
	}
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
	db.loadValues(&cat)
	return cat, nil
}

func (db *ItemDB) LoadEnums() ([]item.Enum, error) {
	enums := make([]item.Enum, 0)
	err := db.Sql.Query(func(row pgx.Rows) error {
		var enum item.Enum
		err := row.Scan(&enum.Id, &enum.Name)
		if err != nil {
			return err
		}
		if enum.Id == 0 {
			return errors.New("no enum found")
		}
		enums = append(enums, enum)
		return nil
	}, SELECT_ALL_ENUM)
	if err != nil {
		return enums, err
	}
	for i := range enums {
		_ = db.loadValues(&enums[i])
	}
	return enums, nil
}

func (db *ItemDB) DeleteEnumWithId(cid int64) error {
	err := db.checkRefs(cid)
	if err != nil {
		return err
	}

	return nil
}

func (db *ItemDB) loadValues(enum *item.Enum) error {
	enum.Values = make([]item.EnumValue, 0)
	err := db.Sql.Query(func(row pgx.Rows) error {
		var prop item.EnumValue
		err := row.Scan(&prop.Name, &prop.Value)
		if err != nil {
			return err
		}
		enum.Values = append(enum.Values, prop)
		return nil
	}, SELECT_ENUM_VALUES_WITH_CID, enum.Id)
	if err != nil {
		return err
	}
	return nil
}
