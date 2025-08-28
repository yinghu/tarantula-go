package persistence

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	DELETE_REGISTRATION_FROM_APP string = "DELETE FROM item_registration WHERE item_id = $1 AND app = $2 AND env = $3"
	SELECT_REGISTRATION_FROM_APP string = "SELECT item_id,scheduling,start_time,close_time,end_time FROM item_registration WHERE app = $1 AND env= $2"
)

func (db *ItemDB) DeleteRegistration(itemId int64, app string, env string) error {
	r, err := db.Sql.Exec(DELETE_REGISTRATION_FROM_APP, itemId, app, env)
	if err != nil {
		return err
	}
	if r == 0 {
		return fmt.Errorf("no row deleted")
	}
	return nil
}
func (db *ItemDB) SaveRegistration(reg item.ConfigRegistration) error {
	if reg.Scheduling {
		_, err := db.Sql.Exec(INSERT_REGISTRATION, reg.ItemId, reg.App, reg.Env, true, reg.StartTime.UnixMilli(), reg.CloseTime.UnixMilli(), reg.EndTime.UnixMilli())
		if err != nil {
			return err
		}
		return nil
	}
	_, err := db.Sql.Exec(INSERT_REGISTRATION, reg.ItemId, reg.App, reg.Env, false, 0, 0, 0)
	if err != nil {
		return err
	}
	return nil
}
func (db *ItemDB) LoadRegistrations(app string, env string) ([]item.ConfigRegistration, error) {
	regs := make([]item.ConfigRegistration, 0)
	err := db.Sql.Query(func(row pgx.Rows) error {
		reg := item.ConfigRegistration{}
		err := row.Scan(&reg.ItemId, &reg.Scheduling, &reg.StartTime, &reg.CloseTime, &reg.EndTime)
		if err != nil {
			return err
		}
		if reg.ItemId == 0 {
			return fmt.Errorf("no item id associated on %s : %s", app, env)
		}
		core.AppLog.Printf("LOADED : %d\n", reg.ItemId)
		regs = append(regs, reg)
		return nil
	}, SELECT_REGISTRATION_FROM_APP, app, env)
	if err != nil {
		return regs, err
	}
	return regs, nil
}
