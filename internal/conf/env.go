package conf

import (
	"encoding/json"
	"io"
	"os"
)

type Env struct {
	GroupName       string   `json:"GroupName"`
	PartitionNumber int      `json:"PartitionNumber"`
	NodeName        string   `json:"NodeName"`
	NodeId          int64    `json:"NodeId"`
	HttpEndpoint    string   `json:"HttpEndpoint"`
	TcpEndpoint     string   `json:"TcpEndpoint"`
	DatabaseURL     string   `json:"DatabaseURL"`
	EtcdEndpoints   []string `json:"EtcdEndpoints"`
}

func (f *Env) Load(fn string) error {
	conf, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer conf.Close()
	data, _ := io.ReadAll(conf)
	json.Unmarshal(data, f)
	if f.PartitionNumber == 0 {
		f.PartitionNumber = 5
	}

	return nil
}
