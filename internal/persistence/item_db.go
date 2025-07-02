package persistence

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
