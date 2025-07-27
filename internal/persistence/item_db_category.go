package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INSERT_CATEGORY                  string = "INSERT INTO item_category (id,name,scope,scope_sequence,rechargeable,description) VALUES($1,$2,$3,$4,$5,$6)"
	INSERT_PROPERTY                  string = "INSERT INTO item_category_property (category_id,name,type,reference,nullable) VALUES($1,$2,$3,$4,$5)"
	SELECT_CATEGORY_WITH_NAME        string = "SELECT id,scope,scope_sequence,rechargeable,description FROM item_category WHERE name = $1"
	SELECT_CATEGORY_WITH_ID          string = "SELECT name,scope,scope_sequence,rechargeable,description FROM item_category WHERE id = $1"
	SELECT_PROPERTIES_WITH_CID       string = "SELECT name,type,reference,nullable FROM item_category_property WHERE category_id = $1"
	SELECT_CATEGORIES_WITH_SCOPE     string = "SELECT id,name,scope,scope_sequence,rechargeable,description FROM item_category WHERE scope_sequence < $1 OR scope = $2"
	DELETE_CATEGORY_WITH_ID          string = "DELETE FROM item_category WHERE id = $1"
	DELETE_CATEGORY_PROPERTY_WITH_ID string = "DELETE FROM item_category_property WHERE category_id = $1"
)

func (db *ItemDB) SaveCategory(c item.Category) error {
	refids, err := db.validateCategory(c)
	if err != nil {
		return err
	}
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
		for i := range refids {
			f, err := tx.Exec(context.Background(), INSERT_REFERENCE, c.Id, refids[i])
			if err != nil {
				return err
			}
			if f.RowsAffected() == 0 {
				return errors.New("no reference insert")
			}
		}
		return db.Gis.SaveCategory(c)
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

func (db *ItemDB) LoadCategories(scopeEnd int32, targetScope string) []item.Category {
	list := make([]item.Category, 0)
	db.Sql.Query(func(row pgx.Rows) error {
		var cat = item.Category{}
		err := row.Scan(&cat.Id, &cat.Name, &cat.Scope, &cat.ScopeSequence, &cat.Rechargeable, &cat.Description)
		if err != nil {
			return err
		}
		list = append(list, cat)
		return nil
	}, SELECT_CATEGORIES_WITH_SCOPE, scopeEnd, targetScope)
	return list
}

func (db *ItemDB) DeleteCategoryWithId(cid int64) error {
	var refs int
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&refs)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_REFERENCE_WITH_REF_ID, cid)
	if err != nil {
		return err
	}
	if refs > 0 {
		return fmt.Errorf("reference ct : %d", refs)
	}
	err = db.Sql.Txn(func(tx pgx.Tx) error {
		dc, err := tx.Exec(context.Background(), DELETE_CATEGORY_WITH_ID, cid)
		if err != nil {
			return err
		}
		if dc.RowsAffected() == 0 {
			return fmt.Errorf("not existed %d", cid)
		}
		pc, err := tx.Exec(context.Background(), DELETE_CATEGORY_PROPERTY_WITH_ID, cid)
		if err != nil {
			return err
		}
		if pc.RowsAffected() == 0 {
			return fmt.Errorf("not property existed %d", cid)
		}
		_, err = tx.Exec(context.Background(), DELETE_REFERENCE_WITH_ITEM_ID, cid)
		if err != nil {
			return err
		}
		err = db.Gis.RemoveCategory(cid)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *ItemDB) validateCategory(c item.Category) ([]int64, error) {
	refids := make([]int64, 0)
	if c.Id <= 0 {
		return refids, errors.New("none negative id required")
	}
	if c.Scope == "" {
		return refids, errors.New("scope none empty string required")
	}
	if c.Name == "" {
		return refids, errors.New("name none empty string required")
	}
	if c.Description == "" {
		return refids, errors.New("description none empty string required")
	}
	if len(c.Properties) == 0 {
		return refids, errors.New("at least 1 property required")
	}
	for i := range c.Properties {
		prop := c.Properties[i]
		if prop.Name == "" {
			return refids, errors.New("prop name none empty string required")
		}
		if prop.Type == "" {
			return refids, errors.New("prop type none empty string required")
		}
		if prop.Reference == "" {
			return refids, errors.New("prop reference none empty string required")
		}
		if prop.Type == "enum" {
			enm, err := db.LoadEnum(prop.Reference)
			if err != nil {
				return refids, err
			}
			refids = append(refids, enm.Id)
			continue
		}
		if prop.Type == "category" {
			parts := strings.Split(prop.Reference, ":")
			if len(parts) != 2 {
				return refids, errors.New("wrong category reference format")
			}
			cat, err := db.LoadCategory(parts[1])
			if err != nil {
				return refids, err
			}
			refids = append(refids, cat.Id)
			continue
		}
		if prop.Type == "set" || prop.Type == "list" {
			parts := strings.Split(prop.Reference, ":")
			if len(parts) != 2 {
				return refids, errors.New("wrong category reference format")
			}
			cat, err := db.LoadCategory(parts[1])
			if err != nil {
				return refids, err
			}
			refids = append(refids, cat.Id)
			continue
		}
		if prop.Type == "scope" {
			parts := strings.Split(prop.Reference, ":")
			if len(parts) != 2 {
				return refids, errors.New("wrong scope reference format")
			}
			continue
		}
	}
	return refids, nil
}
