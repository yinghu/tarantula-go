package persistence

import (
	"fmt"
	"os"

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
	fmt.Printf("Save to %s\n", fn)
	dest, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Err %s\n", err.Error())
		return err
	}
	defer dest.Close()
	dest.WriteString(string(util.ToJson(c)))
	util.GitAdd(fn)
	util.GitCommit("save category [" + fn + "]")
	util.GitPush()
	return nil
}

func (db *GitItemStore) SaveEnum(c item.Enum) error {
	return nil
}

func (db *GitItemStore) SaveConfiguration(c item.Configuration) error {
	return nil
}
