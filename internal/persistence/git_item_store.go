package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type InventoryResp struct {
	core.OnSession
	Stock []item.Inventory
}

type GitItemStore struct {
	RepositoryDir    string
	EnumDir          string
	CategoryDir      string
	ConfigurationDir string
	core.JsonRequester
}

func (db *GitItemStore) Start() error {
	db.CategoryDir = db.RepositoryDir + "/category"
	db.EnumDir = db.RepositoryDir + "/enum"
	db.ConfigurationDir = db.RepositoryDir + "/configuration"
	os.Chdir(db.RepositoryDir)
	util.GitPull()
	return nil
}

func (db *GitItemStore) SaveCategory(c item.Category) error {
	fn := fmt.Sprintf("%s/%d.json", db.CategoryDir, c.Id)
	err := db.writeFile(fn, string(util.ToJson(c)))
	if err != nil {
		core.AppLog.Printf("Err %s\n", err.Error())
		return err
	}
	gr := util.GitAdd(fn)
	if !gr.Successful {
		return errors.New("cannot add file [" + fn + "]")
	}
	idx := fmt.Sprintf("%s/%s.json", db.CategoryDir, c.Name)
	err = db.writeFile(idx, string(util.ToJson(c)))
	if err != nil {
		return err
	}
	gr = util.GitAdd(idx)
	if !gr.Successful {
		return errors.New("cannot add file [" + idx + "]")
	}
	gr = util.GitCommit(fmt.Sprintf("save category [ %s : %d]", c.Name, c.Id))
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	return nil
}

func (db *GitItemStore) SaveEnum(c item.Enum) error {
	fn := fmt.Sprintf("%s/%d.json", db.EnumDir, c.Id)
	dest, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		core.AppLog.Printf("Err %s\n", err.Error())
		return err
	}
	defer dest.Close()
	_, err = dest.WriteString(string(util.ToJson(c)))
	if err != nil {
		return err
	}
	gr := util.GitAdd(fn)
	if !gr.Successful {
		return errors.New("cannot add file [" + fn + "]")
	}
	gr = util.GitCommit(fmt.Sprintf("save enum [ %s : %d]", c.Name, c.Id))
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	return nil
}

func (db *GitItemStore) SaveConfiguration(c item.Configuration) error {
	fn := fmt.Sprintf("%s/%d.json", db.ConfigurationDir, c.Id)
	dest, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		core.AppLog.Printf("Err %s\n", err.Error())
		return err
	}
	defer dest.Close()
	_, err = dest.WriteString(string(util.ToJson(c)))
	if err != nil {
		return err
	}
	gr := util.GitAdd(fn)
	if !gr.Successful {
		return errors.New("cannot add file [" + fn + "]")
	}
	gr = util.GitCommit(fmt.Sprintf("save configuration [ %s : %d]", c.Name, c.Id))
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	return nil
}

func (db *GitItemStore) Load(cid int64) (item.Configuration, error) {
	conf := item.Configuration{Id: cid}
	fn := fmt.Sprintf("%s/%d.json", db.ConfigurationDir, cid)
	src, err := os.Open(fn)
	if err != nil {
		core.AppLog.Printf("Err %s\n", err.Error())
		return conf, err
	}
	defer src.Close()
	err = json.NewDecoder(src).Decode(&conf)
	if err != nil {
		return conf, err
	}
	conf.Reference = map[string]any{}
	for k := range conf.Application {
		refs := conf.Application[k]
		confs := make([]item.Configuration, 0)
		for r := range refs {
			cid, _ := strconv.ParseInt(refs[r], 10, 64)
			conf, err := db.Load(cid)
			if err != nil {
				core.AppLog.Printf("Err %s\n", err.Error())
			}
			confs = append(confs, conf)
		}
		conf.Reference[k] = confs
	}
	return conf, nil
}

func (db *GitItemStore) RemoveCategory(cid int64, cn string) error {
	fn := fmt.Sprintf("%s/%d.json", db.CategoryDir, cid)
	gr := util.GitRemove(fn)
	if !gr.Successful {
		return fmt.Errorf("cannot remove file :%s", fn)
	}
	idx := fmt.Sprintf("%s/%s.json", db.CategoryDir, cn)
	gr = util.GitRemove(idx)
	if !gr.Successful {
		return fmt.Errorf("cannot remove file :%s", idx)
	}
	gr = util.GitCommit(fmt.Sprintf("remove category [%d]", cid))
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	return nil
}
func (db *GitItemStore) RemoveConfig(cid int64) error {
	fn := fmt.Sprintf("%s/%d.json", db.ConfigurationDir, cid)
	gr := util.GitRemove(fn)
	if !gr.Successful {
		return fmt.Errorf("cannot remove file :%s", fn)
	}
	gr = util.GitCommit(fmt.Sprintf("remove configuration [%d]", cid))
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	return nil
}
func (db *GitItemStore) RemoveEnum(cid int64) error {
	fn := fmt.Sprintf("%s/%d.json", db.EnumDir, cid)
	gr := util.GitRemove(fn)
	if !gr.Successful {
		return fmt.Errorf("cannot remove file :%s", fn)
	}
	gr = util.GitCommit(fmt.Sprintf("remove enum [%d]", cid))
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	return nil
}
func (db *GitItemStore) Reload(kv item.KVUpdate) error {
	util.GitPull()
	var repo item.RepoUpdate
	json.Unmarshal([]byte(kv.Value), &repo)
	core.AppLog.Printf("Repo : %v\n", repo)
	return nil
}
func (db *GitItemStore) LoadCategory(name string) (item.Category, error) {
	cat := item.Category{}
	fn := fmt.Sprintf("%s/%s.json", db.CategoryDir, name)
	src, err := os.Open(fn)
	if err != nil {
		return cat, err
	}
	defer src.Close()
	err = json.NewDecoder(src).Decode(&cat)
	if err != nil {
		return cat, err
	}
	return cat, nil
}
func (db *GitItemStore) writeFile(fn string, data string) error {
	dest, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		core.AppLog.Printf("open file error %s\n", err.Error())
		return err
	}
	defer dest.Close()
	_, err = dest.WriteString(data)
	if err != nil {
		core.AppLog.Printf("write file error %s\n", err.Error())
		return err
	}
	return nil
}

func (db *GitItemStore) Grant(inv item.OnInventory) error {
	var er error
	for i := range 5 {
		ret := db.PostJsonSync("http://inventory:8080/inventory/grant", inv)
		if ret.ErrorCode == 0 {
			return nil
		}
		time.Sleep(1000 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
		er = fmt.Errorf("failed on retries %d : %s", i, ret.Message)
	}
	return er
}

func (db *GitItemStore) Validate(c item.Configuration, validator item.Validator) {
	item.ItemValidator(c, validator)
}

func (db *GitItemStore) Stock(inv item.OnInventory) ([]item.Inventory, error) {
	stock := make([]item.Inventory, 0)
	for i := range 5 {
		ret := db.PostJsonSync("http://inventory:8080/inventory/load", inv)
		if ret.ErrorCode == 0 {
			//return nil
		}
		time.Sleep(1000 * time.Millisecond)
		core.AppLog.Printf("Retries: %d %v\n", i, ret)
		//er = fmt.Errorf("failed on retries %d : %s", i, ret.Message)
	}
	return stock, nil
}
