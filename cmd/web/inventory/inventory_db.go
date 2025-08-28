package main

import (
	"context"
	"fmt"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INVENTORY_SQL_SCHEMA      string = "CREATE TABLE IF NOT EXISTS inventory (id SERIAL PRIMARY KEY,system_id BIGINT NOT NULL,type_id VARCHAR(50) NOT NULL, rechargeable BOOL NOT NULL,amount INTEGER NOT NULL,update_time BIGINT NOT NULL, UNIQUE(system_id,type_id))"
	INVENTORY_ITEM_SQL_SCHEMA string = "CREATE TABLE IF NOT EXISTS inventory_item (id SERIAL PRIMARY KEY,inventory_id INTEGER NOT NULL,item_id BIGINT NOT NULL)"

	UPDATE_INVENTORY                string = "INSERT INTO inventory AS iv (system_id,type_id,rechargeable,amount,update_time) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (system_id,type_id) DO UPDATE SET amount = iv.amount + $6 , update_time = $7 WHERE iv.system_id = $8 AND iv.type_id = $9 RETURNING id,amount"
	INSERT_INVENTORY_ITEM           string = "INSERT INTO inventory_item (inventory_id,item_id) VALUES ($1,$2)"
	SELECT_INVENTORY_WITH_SYSTEM_ID string = "SELECT id, type_id, rechargeable,amount,update_time FROM inventory WHERE system_id = $1"
	SELECT_INVENTORY_WITH_TYPE_ID   string = "SELECT id,rechargeable,amount,update_time FROM inventory WHERE system_id = $1 AND type_id = $2"
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

func (s *InventoryService) loadInventory(iv item.OnInventory) ([]item.Inventory, error) {
	ilist := make([]item.Inventory, 0)
	err := s.Sql.Query(func(row pgx.Rows) error {
		var inv item.Inventory
		var t int64
		err := row.Scan(&inv.Id, &inv.TypeId, &inv.Rechargeable, &inv.Amount, &t)
		if err != nil {
			return err
		}
		if inv.Id == 0 {
			return fmt.Errorf("no row id associated")
		}
		inv.UpdateTime = time.UnixMilli(t)
		ilist = append(ilist, inv)
		return nil
	}, SELECT_INVENTORY_WITH_SYSTEM_ID, iv.SystemId)
	if err != nil {
		return ilist, err
	}
	return ilist, nil
}

func (s *InventoryService) updateInventory(iv item.Inventory) error {
	var id int32
	var amount int32
	err := s.Sql.Txn(func(tx pgx.Tx) error {
		err := tx.QueryRow(context.Background(), UPDATE_INVENTORY, iv.SystemId, iv.TypeId, iv.Rechargeable, iv.Amount, iv.UpdateTime.UnixMilli(), iv.Amount, iv.UpdateTime.UnixMilli(), iv.SystemId, iv.TypeId).Scan(&id, &amount)
		if err != nil {
			return err
		}
		if id == 0 {
			return fmt.Errorf("no row updated")
		}
		e, err := tx.Conn().Exec(context.Background(), INSERT_INVENTORY_ITEM, id, iv.ItemId)
		if err != nil {
			return err
		}
		if e.RowsAffected() == 0 {
			return fmt.Errorf("no row updated")
		}
		return nil
	})
	if err != nil {
		return err
	}
	core.AppLog.Printf("Id %d\n", id)
	return nil
}
