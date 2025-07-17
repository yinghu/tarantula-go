package persistence

import (
	"context"
	"errors"
	"strings"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INSERT_CATEGORY              string = "INSERT INTO item_category (id,name,scope,scope_sequence,rechargeable,description) VALUES($1,$2,$3,$4,$5,$6)"
	INSERT_PROPERTY              string = "INSERT INTO item_category_property (category_id,name,type,reference,nullable) VALUES($1,$2,$3,$4,$5)"
	SELECT_CATEGORY_WITH_NAME    string = "SELECT id,scope,scope_sequence,rechargeable,description FROM item_category WHERE name = $1"
	SELECT_CATEGORY_WITH_ID      string = "SELECT name,scope,scope_sequence,rechargeable,description FROM item_category WHERE id = $1"
	SELECT_PROPERTIES_WITH_CID   string = "SELECT name,type,reference,nullable FROM item_category_property WHERE category_id = $1"
	SELECT_CATEGORIES_WITH_SCOPE string = "SELECT id,name,scope,scope_sequence,rechargeable,description FROM item_category WHERE scope_sequence >= $1 AND scope_sequence <= $2"
)

func (db *ItemDB) SaveCategory(c item.Category) error {
	return db.Sql.Txn(func(tx pgx.Tx) error {
		r, err := tx.Exec(context.Background(), INSERT_CATEGORY, c.Id, c.Name, c.Scope, c.ScopeSequence, c.Rechargeable, c.Description)
		if err != nil {
			return err
		}
		if r.RowsAffected() == 0 {
			return errors.New("no insert")
		}
		valid := false
		for i := range c.Properties {
			p := c.Properties[i]
			_, err = tx.Exec(context.Background(), INSERT_PROPERTY, c.Id, p.Name, p.Type, p.Reference, p.Nullable)
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

func (db *ItemDB) LoadCategoryWithId(cid int64) (item.Category, error) {
	cat := item.Category{Id: cid}
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&cat.Name, &cat.Scope, &cat.ScopeSequence, &cat.Rechargeable, &cat.Description)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_CATEGORY_WITH_ID, cid)
	if err != nil {
		return cat, err
	}
	if cat.Name == "" {
		return cat, errors.New("category not existed")
	}
	cat.Properties = make([]item.Property, 0)
	err = db.Sql.Query(func(row pgx.Rows) error {
		var prop item.Property
		err := row.Scan(&prop.Name, &prop.Type, &prop.Reference, &prop.Nullable)
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

func (db *ItemDB) LoadCategory(cname string) (item.Category, error) {
	cat := item.Category{Name: cname}
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&cat.Id, &cat.Scope, &cat.ScopeSequence, &cat.Rechargeable, &cat.Description)
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
		err := row.Scan(&prop.Name, &prop.Type, &prop.Reference, &prop.Nullable)
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

func (db *ItemDB) LoadCategories(scopeStart int32, scopeEnd int32) []item.Category {
	list := make([]item.Category, 0)
	db.Sql.Query(func(row pgx.Rows) error {
		var cat = item.Category{}
		err := row.Scan(&cat.Id, &cat.Name, &cat.Scope, &cat.ScopeSequence, &cat.Rechargeable, &cat.Description)
		if err != nil {
			return err
		}
		list = append(list, cat)
		return nil
	}, SELECT_CATEGORIES_WITH_SCOPE, scopeStart, scopeEnd)
	return list
}

func (db *ItemDB) ValidateCategory(c item.Category) error {
	if c.Id <= 0 {
		return errors.New("none negative id required")
	}
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
		if prop.Type == "scope" {
			parts := strings.Split(prop.Reference, ":")
			if len(parts) != 2 {
				return errors.New("wrong scope reference format")
			}
			continue
		}
	}
	return nil
}
