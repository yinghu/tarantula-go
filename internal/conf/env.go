package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gameclustering.com/internal/core"
)

const (
	NODE_GROUP string = "TN_GROUP"
	NODE_HOST  string = "TN_HOST"
	NODE_NAME  string = "TN_NAME"
	NODE_ID    string = "TN_ID"
)

type Sql struct {
	DatabaseURL string `json:"DatabaseURL"`
}

type LocalStore struct {
	InMemory bool `json:"InMemory"`
}

type EventEndpoint struct {
	Enabled         bool   `json:"Enabled"`
	OutboundEnabled bool   `json:"OutboundEnabled"`
	TcpEndpoint     string `json:"TcpEndpoint"`
}

type Env struct {
	Prefix        string        `json:"Prefix"`
	Clustering    bool          `json:"Clustering"`
	Standalone    bool          `json:"Standalone"`
	GroupName     string        `json:"GroupName"`
	NodeName      string        `json:"NodeName"`
	NodeId        int64         `json:"NodeId"`
	LocalDir      string        `json:"LocalDir"`
	HttpBinding   string        `json:"HttpBinding"`
	HttpEndpoint  string        `json:"HttpEndpoint"`
	Evp           EventEndpoint `json:"EventEndpoint"`
	EtcdEndpoints []string      `json:"EtcdEndpoints"`
	ManagedApps   []string      `json:"ManagedApps"`
	Pgs           Sql           `json:"Sql"`
	Bdg           LocalStore    `json:"LocalStore"`
	Bin           string        `json:"-"`
	HomeDir       string        `json:"-"`
}

func (f *Env) ClusterCtx() string {
	return f.Prefix + "/" + f.GroupName
}

func (f *Env) PresenceCtx() string {
	return f.Prefix + "/presence"
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
		f.HomeDir = homeDir
		err = os.MkdirAll(f.LocalDir, 0755)
		if err != nil {
			return err
		}
	}
	h, eh := os.LookupEnv(NODE_HOST)
	if eh {
		fmt.Printf("Using http endpoint : %s\n", h)
		f.HttpEndpoint = h
		parts := strings.Split(f.Evp.TcpEndpoint, ":")
		f.Evp.TcpEndpoint = parts[0] + "://" + h + ":" + parts[2]
		fmt.Printf("Using tcp endpoint : %s\n", f.Evp.TcpEndpoint)
	}
	n, en := os.LookupEnv(NODE_NAME)
	if en {
		fmt.Printf("Using node name : %s\n", n)
		f.NodeName = n
	}
	d, ed := os.LookupEnv(NODE_ID)
	if ed {
		fmt.Printf("Using node id : %s\n", d)
		id, err := strconv.Atoi(d)
		if err == nil {
			f.NodeId = int64(id)
			fmt.Printf("Node id : %d %d\n", id, f.NodeId)
		}
	}
	g, eg := os.LookupEnv(NODE_GROUP)
	if eg {
		fmt.Printf("Using node group prefix : %s\n", g)
		f.Prefix = g
	}
	if f.Prefix == "" {
		f.Prefix = "dev"
	}
	core.CreateAppLog(f.LocalDir)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	f.Bin = wd
	return nil
}
