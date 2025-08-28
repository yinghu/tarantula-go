package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"github.com/jackc/pgx/v5"
)

const (
	INSERT_CONFIG       string = "INSERT INTO item_configuration (id,name,type,type_id,category,version) VALUES($1,$2,$3,$4,$5,$6)"
	INSERT_HEADER       string = "INSERT INTO item_header (configuration_id,name,value) VALUES($1,$2,$3)"
	INSERT_APPLICATION  string = "INSERT INTO item_application (configuration_id,name,reference_id) VALUES($1,$2,$3)"
	INSERT_REGISTRATION string = "INSERT INTO item_registration (item_id,app,env,scheduling,start_time,close_time,end_time) VALUES($1,$2,$3,$4,$5,$6,$7)"

	SELECT_CONFIG_WITH_NAME           string = "SELECT id,name,type,type_id,version FROM item_configuration WHERE category = $1 LIMIT $2"
	SELECT_CONFIG_WITH_ID             string = "SELECT name,type,type_id,category,version FROM item_configuration WHERE id = $1"
	SELECT_CONFIG_HEADER_WIHT_ID      string = "SELECT name,value FROM item_header WHERE configuration_id = $1"
	SELECT_CONFIG_APPLICATION_WITH_ID string = "SELECT name,reference_id FROM item_application WHERE configuration_id = $1"

	DELETE_HEADER         string = "DELETE FROM item_header WHERE configuration_id = $1"
	DELETE_APPLICATION    string = "DELETE FROM item_application WHERE configuration_id = $1"
	DELETE_CONFIG_WITH_ID string = "DELETE FROM item_configuration WHERE id = $1"

	SELECT_REGISTRATION_WITH_ITEM_ID_APP string = "SELECT id,scheduling,start_time,close_time,end_time FROM item_registration WHERE item_id = $1 AND app = $2 AND env= $3"
	SELECT_REGISTRATION_WITH_ITEM_ID     string = "SELECT COUNT(*) FROM item_registration WHERE item_id = $1"
	DELETE_REGISTRATION_WITH_ID          string = "DELETE FROM item_registration AS d WHERE id = $1 RETURNING d.item_id, d.app, d.env"
	DELETE_REGISTRATION_FROM_APP         string = "DELETE FROM item_registration WHERE item_id = $1 AND app = $2 AND env = $3 AND"
)

func (db *ItemDB) Save(c item.Configuration) error {
	refids, err := db.validate(c)
	if err != nil {
		return err
	}
	return db.Sql.Txn(func(tx pgx.Tx) error {
		r, err := tx.Exec(context.Background(), INSERT_CONFIG, c.Id, c.Name, c.Type, c.TypeId, c.Category, c.Version)
		if err != nil {
			return err
		}
		if r.RowsAffected() == 0 {
			return errors.New("no insert")
		}
		for k, v := range c.Header {
			inserted, err := tx.Exec(context.Background(), INSERT_HEADER, c.Id, k, fmt.Sprintf("%v", v))
			if err != nil {
				return err
			}
			if inserted.RowsAffected() != 1 {
				return errors.New("no data inserted")
			}
		}
		for k, v := range c.Application {
			for i := range v {
				aid, err := strconv.ParseInt(v[i], 10, 64)
				if err != nil {
					return err
				}
				inserted, err := tx.Exec(context.Background(), INSERT_APPLICATION, c.Id, k, aid)
				if err != nil {
					return err
				}
				if inserted.RowsAffected() != 1 {
					return errors.New("no data inserted")
				}
			}
		}
		for i := range refids {
			//fmt.Printf("ref id : %d=>%d\n", i, refids[i])
			f, err := tx.Exec(context.Background(), INSERT_REFERENCE, c.Id, refids[i])
			if err != nil {
				return err
			}
			if f.RowsAffected() == 0 {
				return errors.New("no reference insert")
			}
		}
		return db.Gis.SaveConfiguration(c)
	})
}

func (db *ItemDB) LoadWithName(cname string, limit int) ([]item.Configuration, error) {
	list := make([]item.Configuration, limit)
	ct := 0
	err := db.Sql.Query(func(row pgx.Rows) error {
		conf := item.Configuration{Category: cname}
		err := row.Scan(&conf.Id, &conf.Name, &conf.Type, &conf.TypeId, &conf.Version)
		if err != nil {
			return err
		}
		if conf.Id > 0 {
			list[ct] = conf
			ct++
		}
		return nil
	}, SELECT_CONFIG_WITH_NAME, cname, limit)
	if err != nil {
		return nil, err
	}
	for i := range list[:ct] {
		db.loadHeader(&list[i])
		db.loadApplication(&list[i])
		list[i].Reference = map[string]any{}
		for k := range list[i].Application {
			refs := list[i].Application[k]
			confs := make([]item.Configuration, 0)
			for r := range refs {
				//fmt.Printf("Ref 1 id : %s\n", refs[r])
				cid, _ := strconv.ParseInt(refs[r], 10, 64)
				conf, err := db.LoadWithId(cid)
				if err != nil {
					core.AppLog.Printf("Err %s\n", err.Error())
				}
				confs = append(confs, conf)
			}
			list[i].Reference[k] = confs
		}
	}
	return list[:ct], nil

}

func (db *ItemDB) LoadWithId(cid int64) (item.Configuration, error) {
	conf := item.Configuration{Id: cid}
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&conf.Name, &conf.Type, &conf.TypeId, &conf.Category, &conf.Version)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_CONFIG_WITH_ID, cid)
	if err != nil {
		return conf, err
	}
	if conf.Name == "" {
		return conf, errors.New("obj not existed")
	}
	db.loadHeader(&conf)
	db.loadApplication(&conf)
	conf.Reference = map[string]any{}
	for k := range conf.Application {
		refs := conf.Application[k]
		confs := make([]item.Configuration, 0)
		for r := range refs {
			cid, _ := strconv.ParseInt(refs[r], 10, 64)
			conf, err := db.LoadWithId(cid)
			if err != nil {
				core.AppLog.Printf("Err %s\n", err.Error())
			}
			confs = append(confs, conf)
		}
		conf.Reference[k] = confs
	}
	return conf, nil
}

func (db *ItemDB) DeleteWithId(cid int64) error {
	err := db.checkRefs(cid)
	if err != nil {
		return err
	}
	err = db.checkRegs(cid)
	if err != nil {
		return err
	}
	return db.Sql.Txn(func(tx pgx.Tx) error {
		_, err := tx.Exec(context.Background(), DELETE_CONFIG_WITH_ID, cid)
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), DELETE_HEADER, cid)
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), DELETE_APPLICATION, cid)
		if err != nil {
			return err
		}
		_, err = tx.Exec(context.Background(), DELETE_REFERENCE_WITH_ITEM_ID, cid)
		if err != nil {
			return err
		}
		err = db.Gis.RemoveConfig(cid)
		if err != nil {
			return err
		}
		return nil
	})
}

func (db *ItemDB) Register(reg item.ConfigRegistration) error {
	conf, err := db.LoadWithId(reg.ItemId)
	if err != nil {
		return err
	}
	sc, ok := conf.Reference["Schedule"].([]item.Configuration)
	reg.Scheduling = ok
	if reg.Scheduling {
		jsc, err := json.Marshal(sc[0].Header)
		if err != nil {
			core.AppLog.Printf("no schedule data\n")
			reg.Scheduling = false
		} else {
			err = json.Unmarshal(jsc, &reg)
			if err != nil {
				core.AppLog.Printf("no schedule data %s\n", err.Error())
				reg.Scheduling = false
			}
		}
	}
	if reg.Scheduling {
		_, err := db.Sql.Exec(INSERT_REGISTRATION, reg.ItemId, reg.App, reg.Env, true, reg.StartTime.UnixMilli(), reg.CloseTime.UnixMilli(), reg.EndTime.UnixMilli())
		if err != nil {
			return err
		}
		db.Schedule(reg)
		return nil
	}
	_, err = db.Sql.Exec(INSERT_REGISTRATION, reg.ItemId, reg.App, reg.Env, false, 0, 0, 0)
	if err != nil {
		return err
	}
	db.Schedule(reg)
	return nil
}

func (db *ItemDB) checkRegs(itemId int64) error {
	var ct int32
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&ct)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_REGISTRATION_WITH_ITEM_ID, itemId)
	if err != nil {
		return err
	}
	if ct > 0 {
		return fmt.Errorf("register ct %d", ct)
	}
	return nil
}

func (db *ItemDB) Check(reg item.ConfigRegistration) (item.ConfigRegistration, error) {
	//reg := item.ConfigRegistration{ItemId: itemId, App: app}
	err := db.Sql.Query(func(row pgx.Rows) error {
		err := row.Scan(&reg.Id, &reg.Scheduling, &reg.StartTime, &reg.CloseTime, &reg.EndTime)
		if err != nil {
			return err
		}
		return nil
	}, SELECT_REGISTRATION_WITH_ITEM_ID_APP, reg.ItemId, reg.App, reg.Env)
	if err != nil {
		return reg, err
	}
	if reg.Id == 0 {
		return reg, errors.New("obj not existed")
	}
	return reg, nil
}
func (db *ItemDB) Release(regId int32) error {
	var deleted item.ConfigRegistration
	err := db.Sql.Txn(func(tx pgx.Tx) error {
		return tx.QueryRow(context.Background(), DELETE_REGISTRATION_WITH_ID, regId).Scan(&deleted.ItemId, &deleted.App, &deleted.Env)
	})
	if err != nil {
		return err
	}
	if deleted.ItemId == 0 {
		return errors.New("no row deleted")
	}
	db.Unschedule(deleted)
	return nil
}

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

func (db *ItemDB) loadHeader(c *item.Configuration) error {
	c.Header = make(map[string]any)
	return db.Sql.Query(func(row pgx.Rows) error {
		var k string
		var v any
		err := row.Scan(&k, &v)
		if err != nil {
			return err
		}
		c.Header[k] = v
		return nil
	}, SELECT_CONFIG_HEADER_WIHT_ID, c.Id)
}

func (db *ItemDB) loadApplication(c *item.Configuration) error {
	c.Application = make(map[string][]string)
	return db.Sql.Query(func(row pgx.Rows) error {
		var k string
		var v int64
		err := row.Scan(&k, &v)
		if err != nil {
			return err
		}
		c.Application[k] = append(c.Application[k], fmt.Sprintf("%d", v))
		return nil
	}, SELECT_CONFIG_APPLICATION_WITH_ID, c.Id)
}

func (db *ItemDB) validate(c item.Configuration) ([]int64, error) {
	refids := make([]int64, 0)
	if c.Name == "" {
		return refids, errors.New("name none empty string required")
	}
	if c.TypeId == "" {
		return refids, errors.New("typeId none empty string required")
	}
	if c.Type == "" {
		return refids, errors.New("type none empty string required")
	}
	if c.Category == "" {
		return refids, errors.New("category none empty string required")
	}
	if c.Version == "" {
		return refids, errors.New("version none empty string required")
	}
	cat, err := db.LoadCategory(c.Category)
	if err != nil {
		return refids, err
	}
	refids = append(refids, cat.Id)
	valid := len(c.Header)
	for i := range cat.Properties {
		prop := cat.Properties[i]
		if prop.Type == "category" || prop.Type == "set" || prop.Type == "list" || prop.Type == "scope" {
			for _, v := range c.Application {
				for i := range v {
					aid, err := strconv.ParseInt(v[i], 10, 64)
					if err != nil {
						return refids, err
					}
					_, err = db.LoadWithId(aid)
					if err != nil {
						return refids, err
					}
					refids = append(refids, aid)
				}
			}
			continue
		}
		v, existed := c.Header[prop.Name]
		if !existed && !prop.Nullable {
			return refids, errors.New("value not existed : " + prop.Type)
		}
		valid--
		if prop.Type == "string" {
			err = asString(v)
			if err != nil {
				return refids, err
			}
			continue
		}
		if prop.Type == "int" {
			err = asInt(v)
			if err != nil {
				return refids, err
			}
			continue
		}
		if prop.Type == "long" {
			err = asLong(v)
			if err != nil {
				return refids, err
			}
			continue
		}
		if prop.Type == "float" {
			err = asFloat(v)
			if err != nil {
				return refids, err
			}
			continue
		}
		if prop.Type == "double" {
			err = asDouble(v)
			if err != nil {
				return refids, err
			}
			continue
		}
		if prop.Type == "boolean" {
			err = asBool(v)
			if err != nil {
				return refids, err
			}
			continue
		}
		if prop.Type == "dateTime" {
			err = asString(v)
			if err != nil {
				return refids, err
			}
			_, err = time.Parse(time.RFC3339, fmt.Sprintf("%v", v))
			if err != nil {
				return refids, err
			}
			continue
		}
		if prop.Type == "enum" {
			em, err := db.LoadEnum(prop.Reference)
			if err != nil {
				return refids, err
			}
			e, err := toInt32(v)
			if err != nil {
				return refids, err
			}
			matched := false
			for i := range em.Values {
				matched = em.Values[i].Value == e
				if matched {
					break
				}
			}
			if !matched {
				return refids, errors.New("enum value not matched")
			}
			refids = append(refids, em.Id)
			continue
		}
	}
	if valid == 0 {
		slices.Sort(refids)
		refids = slices.Compact(refids)
		return refids, nil
	}
	return refids, errors.New("invalid data")
}

func asString(v any) error {
	_, ok := v.(string)
	if ok {
		return nil
	}
	return errors.New("wrong string format")
}

func asDouble(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong double format")
}

func asFloat(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong float format")
}

func asInt(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong int format")
}

func toInt32(v any) (int32, error) {
	x, err := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 32)
	if err != nil {
		return 0, err

	}
	return int32(x), nil
}

func asLong(v any) error {
	_, ok := v.(float64)
	if ok {
		return nil
	}
	return errors.New("wrong long format")
}

func asBool(v any) error {
	_, ok := v.(bool)
	if ok {
		return nil
	}
	return errors.New("wrong bool format")
}
