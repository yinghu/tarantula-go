package main

import (
	//"errors"
	//"github.com/jackc/pgx/v5"
)

const (
	ITEM_INDEX_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS item_index (id SERIAL PRIMARY KEY,name VARCHAR(100) NOT NULL UNIQUE,category VARCHAR(50) NOT NULL,file_index VARCHAR(100) NOT NULL,update_time TIMESTAMP DEFAULT NOW())"
	INVENTORY_SQL_SCHEMA  string = "CREATE TABLE IF NOT EXISTS inventory (id SERIAL PRIMARY KEY,system_id BIGINT NOT NULL,item_id INTEGER NOT NULL, rechargeable BOOL NOT NULL,update_time TIMESTAMP NOT NULL, UNIQUE(system_id,item_id))"
	//ASSET_INDEX_INSERT_UPDATE string = "INSERT INTO asset_index AS a (system_id,name,file_index,upload_time) VALUES($1,$2,$3,$4) ON CONFLICT (system_id,name) DO UPDATE SET file_index = $5, upload_time = $6  WHERE a.system_id = $7 AND a.name = $8"
	//ASSET_INDEX_SELECT        string = "SELECT file_index FROM asset_index AS a WHERE a.system_id = $1 AND a.name = $2"
)

func (s *InventoryService) createSchema() error {
	
	_, err := s.Sql.Exec(ITEM_INDEX_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(INVENTORY_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil

}
