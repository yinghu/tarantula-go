package persistence

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/item"
	"gameclustering.com/internal/util"
)

type GitItemStore struct {
	RepositoryDir    string
	EnumDir          string
	CategoryDir      string
	ConfigurationDir string
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
	idx := fmt.Sprintf("%s/%s.index", db.CategoryDir, c.Name)
	err = db.writeFile(idx, fmt.Sprintf("%d", c.Id))
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
	idx := fmt.Sprintf("%s/%s.index", db.CategoryDir, cn)
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
	idx := fmt.Sprintf("%s/%s.index", db.CategoryDir, name)
	id, err := db.readIndex(idx)
	if err != nil {
		return cat, err
	}
	fn := fmt.Sprintf("%s/%d.json", db.CategoryDir, id)
	db.readJson(fn, cat)
	core.AppLog.Printf("category %v\n", cat)
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
func (db *GitItemStore) readIndex(fn string) (int64, error) {
	src, err := os.Open(fn)
	if err != nil {
		core.AppLog.Printf("Err %s\n", err.Error())
		return 0, err
	}
	defer src.Close()
	scanner := bufio.NewScanner(src)
	scanner.Scan()
	cid, err := strconv.ParseInt(strings.Split(scanner.Text()," ")[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return cid, nil
}

func (db *GitItemStore) readJson(fn string, t any) error {
	src, err := os.Open(fn)
	if err != nil {
		core.AppLog.Printf("Err %s\n", err.Error())
		return err
	}
	defer src.Close()
	err = json.NewDecoder(src).Decode(&t)
	if err != nil {
		return err
	}
	core.AppLog.Printf("category %v\n", t)
	return nil
}
