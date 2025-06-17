package persistence

import (
	"errors"
	"fmt"

	"gameclustering.com/internal/item"
)

type ItemDB struct {
	Sql *Postgresql
}

func (db *ItemDB) Save(c item.Configuration) error {
	inserted, err := db.Sql.Exec("INSERT INTO item_configuration (name,type,type_id,category,version) VALUES($1,$2,$3,$4,$5) RETURNING id", c.Name, c.Type, c.TypeId, c.Category, c.Version)
	if err != nil {
		return err
	}
	if inserted == 0 {
		return errors.New("config cannot be saved")
	}
	fmt.Printf("id : %d\n",inserted)
	return nil
}
func (db *ItemDB) Load(c item.Configuration) error {
	return nil
}
