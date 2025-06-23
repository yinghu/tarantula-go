package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gameclustering.com/internal/core"
)

const (
	TN_NODE_CONFIG string = "TARANTULA_NODE_CONFIG"
)

type Sql struct {
	DatabaseURL string `json:"DatabaseURL"`
}

type LocalStore struct {
	InMemory bool `json:"InMemory"`
}

type EventEndpoint struct {
	Enabled     bool   `json:"Enabled"`
	TcpEndpoint string `json:"TcpEndpoint"`
}

type Env struct {
	GroupName     string        `json:"GroupName"`
	NodeName      string        `json:"NodeName"`
	NodeId        int64         `json:"NodeId"`
	Presence      string        `json:"Presence"`
	LocalDir      string        `json:"LocalDir"`
	HttpBinding   string        `json:"HttpBinding"`
	HttpEndpoint  string        `json:"HttpEndpoint"`
	Evp           EventEndpoint `json:"EventEndpoint"`
	EtcdEndpoints []string      `json:"EtcdEndpoints"`
	Pgs           Sql           `json:"Sql"`
	Bdg           LocalStore    `json:"LocalStore"`
}

func (f *Env) Load(fn string) error {
	conf, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer conf.Close()
	data, _ := io.ReadAll(conf)
	json.Unmarshal(data, f)
	if f.HttpBinding == "" {
		f.HttpBinding = ":8080"
	}
	if f.LocalDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		f.LocalDir = homeDir + "/" + f.GroupName
		err = os.MkdirAll(f.LocalDir, 0755)
		if err != nil {
			return err
		}
	}
	cfg, existed := os.LookupEnv(TN_NODE_CONFIG)
	if existed {
		parts := strings.Split(cfg, "#")
		fmt.Printf("%s\n", parts[0])
		fmt.Printf("%s\n", parts[1])
		fmt.Printf("%s\n", parts[2])
	}
	core.CreateAppLog(f.LocalDir)
	return nil
}
