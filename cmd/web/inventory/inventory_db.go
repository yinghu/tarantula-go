package main

//"errors"
//"github.com/jackc/pgx/v5"

const (
	INVENTORY_SQL_SCHEMA      string = "CREATE TABLE IF NOT EXISTS inventory (id SERIAL PRIMARY KEY,system_id BIGINT NOT NULL,type_id VARCHAR(50) NOT NULL, rechargeable BOOL NOT NULL,amount INTEGER NOT NULL,update_time TIMESTAMP NOT NULL, UNIQUE(system_id,type_id))"
	INVENTORY_ITEM_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS inventory_item (id SERIAL PRIMARY KEY,inventory_id INTEGER NOT NULL,item_id BIGINT NOT NULL)"
)

func (s *InventoryService) createSchema() error {
	_, err := s.Sql.Exec(INVENTORY_SQL_SCHEMA)
	if err != nil {
		return err
	}
	_, err = s.Sql.Exec(INVENTORY_ITEM_SQL_SCHEMA)
	if err != nil {
		return err
	}
	return nil

}
