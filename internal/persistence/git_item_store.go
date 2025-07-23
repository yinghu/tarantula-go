package persistence

import (
	"errors"
	"fmt"
	"os"

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
	gr = util.GitCommit("save category [" + fn + "]")
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	gr = util.GitPush()
	if !gr.Successful {
		return errors.New("cannot push file [" + fn + "]")
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
	gr = util.GitCommit("save enum [" + fn + "]")
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	gr = util.GitPush()
	if !gr.Successful {
		return errors.New("cannot push file [" + fn + "]")
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
	gr = util.GitCommit("save configuration [" + fn + "]")
	if !gr.Successful {
		return errors.New("cannot commit file [" + fn + "]")
	}
	gr = util.GitPush()
	if !gr.Successful {
		return errors.New("cannot push file [" + fn + "]")
	}
	return nil
}
