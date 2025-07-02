package persistence

import (
	"context"
	"errors"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

func (db *ItemDB) SaveEnum(c item.Enum) error {
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
		return nil
	})
}

func (db *ItemDB) ValidateEnum(c item.Enum) error {
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
