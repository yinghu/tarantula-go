package conf

import (
	"encoding/json"
	"io"
	"os"
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
	return nil
}
