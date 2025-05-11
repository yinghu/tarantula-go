package conf

import (
	"encoding/json"
	"io"
	"os"
)

type Env struct {
	GroupName     string   `json:"GroupName"`
	NodeName      string   `json:"NodeName"`
	NodeId        int64    `json:"NodeId"`
	HttpEndpoint  string   `json:"HttpEndpoint"`
	TcpEndpoint   string   `json:"TcpEndpoint"`
	DatabaseURL   string   `json:"DatabaseURL"`
	EtcdEndpoints []string `json:"EtcdEndpoints"`
}

func (f *Env) Load(fn string) error {
	conf, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer conf.Close()
	data, _ := io.ReadAll(conf)
	json.Unmarshal(data,f)
	//f.GroupName = "presence"
	//f.NodeName = "a01"
	//f.NodeId = 1
	//f.HttpEndpoint = "192.168.1.4:8080"
	//f.DatabaseURL = "postgres://postgres:password@192.168.1.7:5432/tarantula_user"
	//f.EtcdEndpoints = []string{"192.168.1.7:2379"}
	//f.TcpEndpoint = "tcp://192.168.1.4:5000"
	return nil
}
