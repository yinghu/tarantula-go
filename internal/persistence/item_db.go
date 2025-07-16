package persistence

const (
	ITEM_ENUM_SQL_SCHEMA       string = "CREATE TABLE IF NOT EXISTS item_enum (id BIGINT NOT NULL,name VARCHAR(100) NOT NULL UNIQUE,PRIMARY KEY(id))"
	ITEM_ENUM_VALUE_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_enum_value (enum_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,value INTEGER NOT NULL,PRIMARY KEY(enum_id,name))"
	ITEM_CATEGORY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_category (id BIGINT NOT NULL,name VARCHAR(100) NOT NULL UNIQUE,scope VARCHAR(30) NOT NULL ,rechargeable BOOL NOT NULL ,description VARCHAR(100) NOT NULL,PRIMARY KEY(id))"
	ITEM_PROPERTY_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_category_property (category_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,type VARCHAR(100) NOT NULL ,reference VARCHAR(100) NOT NULL ,nullable BOOL NOT NULL ,downloadable BOOL NOT NULL, PRIMARY KEY(category_id,name))"

	ITEM_CONFIGURATION_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_configuration (id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,type VARCHAR(50) NOT NULL ,type_id VARCHAR(50) NOT NULL ,category VARCHAR(100) NOT NULL ,version VARCHAR(10) NOT NULL,UNIQUE(name,version),PRIMARY KEY(id))"
	ITEM_HEADER_SQL_SCHEMA        string = "CREATE TABLE IF NOT EXISTS item_header (configuration_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,value VARCHAR(100) NOT NULL, PRIMARY KEY(configuration_id,name))"
	ITEM_APPLICATION_SQL_SCHEMA   string = "CREATE TABLE IF NOT EXISTS item_application (configuration_id BIGINT NOT NULL,name VARCHAR(100) NOT NULL,reference_id BIGINT NOT NULL,PRIMARY KEY(configuration_id,name,reference_id))"



	DELETE_CONFIG_WITH_NAME string = "DELETE FROM item_configuration WHERE name = $1 RETURNING id"
	DELETE_HEADER           string = "DELETE FROM item_header WHERE configuration_id = $1"
	DELETE_APPLICATION      string = "DELETE FROM item_application WHERE configuration_id = $1"
	DELETE_CONFIG_WITH_ID   string = "DELETE FROM item_configuration WHERE id"
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
