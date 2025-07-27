package persistence

import (
	"gameclustering.com/internal/item"
)

const (
	ITEM_ENUM_SQL_SCHEMA       string = "CREATE TABLE IF NOT EXISTS item_enum (id BIGINT NOT NULL,name VARCHAR(100) NOT NULL UNIQUE,PRIMARY KEY(id))"
	ITEM_ENUM_VALUE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_enum_value (enum_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,value INTEGER NOT NULL,PRIMARY KEY(enum_id,name))"
	ITEM_CATEGORY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_category (id BIGINT NOT NULL,name VARCHAR(100) NOT NULL UNIQUE,scope VARCHAR(30) NOT NULL,scope_sequence INTEGER DEFAULT 0,rechargeable BOOL NOT NULL,description VARCHAR(100) NOT NULL,PRIMARY KEY(id))"
	ITEM_PROPERTY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_category_property (category_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,type VARCHAR(100) NOT NULL ,reference VARCHAR(100) NOT NULL ,nullable BOOL NOT NULL ,PRIMARY KEY(category_id,name))"

	ITEM_CONFIGURATION_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_configuration (id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,type VARCHAR(50) NOT NULL ,type_id VARCHAR(50) NOT NULL ,category VARCHAR(100) NOT NULL ,version VARCHAR(10) NOT NULL,UNIQUE(name,version),PRIMARY KEY(id))"
	ITEM_HEADER_SQL_SCHEMA        string = "CREATE TABLE IF NOT EXISTS item_header (configuration_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,value VARCHAR(100) NOT NULL, PRIMARY KEY(configuration_id,name))"
	ITEM_APPLICATION_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_application (configuration_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,reference_id BIGINT NOT NULL,PRIMARY KEY(configuration_id,name,reference_id))"
	ITEM_REFERENCE_SQL_SCHEMA     string = "CREATE TABLE IF NOT EXISTS item_reference (id SERIAL PRIMARY KEY, item_id BIGINT NOT NULL,ref_id BIGINT NOT NULL)"

	INSERT_REFERENCE              string = "INSERT INTO item_reference (item_id,ref_id) VALUES ($1,$2)"
	SELECT_REFERENCE_WITH_REF_ID  string = "SELECT COUNT(*) FROM item_reference WHERE ref_id = $1"
	DELETE_REFERENCE_WITH_ITEM_ID string = "DELETE FROM item_reference WHERE item_id = $1"
)

type ItemDB struct {
	Sql *Postgresql
	Gis *GitItemStore
}

func (db *ItemDB) Loader() item.ItemLoader {
	return db.Gis
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
	_, err = db.Sql.Exec(ITEM_REFERENCE_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil
}
