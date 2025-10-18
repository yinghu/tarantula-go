package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/util"
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
	LogTruncated  bool          `json:"LogTruncated"`
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
	core.CreateAppLog(f.LocalDir, f.LogTruncated)
	c, exists := os.LookupEnv(NODE_HOST)
	if exists {
		core.AppLog.Printf("Using http endpoint : %s\n", c)
		f.HttpEndpoint = c
		parts := strings.Split(f.Evp.TcpEndpoint, ":")
		f.Evp.TcpEndpoint = parts[0] + "://" + c + ":" + parts[2]
		core.AppLog.Printf("Using tcp endpoint : %s\n", f.Evp.TcpEndpoint)
	}
	c, exists = os.LookupEnv(NODE_NAME)
	if exists {
		core.AppLog.Printf("Using node name : %s\n", c)
		f.NodeName = c
	}
	c, exists = os.LookupEnv(NODE_ID)
	if exists {
		core.AppLog.Printf("Using node id : %s\n", c)
		id, err := strconv.Atoi(c)
		if err == nil {
			f.NodeId = int64(id)
			core.AppLog.Printf("Node id : %d %d\n", id, f.NodeId)
		}
	}
	c, exists = os.LookupEnv(NODE_GROUP)
	if exists {
		core.AppLog.Printf("Using node group prefix : %s\n", c)
		f.Prefix = c
	}
	if f.Prefix == "" {
		f.Prefix = "dev"
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	f.Bin = wd
	//overriding configs
	c, exists = os.LookupEnv("ENV")
	if exists {
		core.AppLog.Printf("Overriding ENV : %s\n", c)
		f.Prefix = c
	}
	c, exists = os.LookupEnv("SEQ")
	if exists {
		core.AppLog.Printf("Overriding SEQ : %s\n", c)
		f.NodeName = fmt.Sprintf("%s.%s", f.GroupName, c)
	}
	c, exists = os.LookupEnv("ETCD_ENDPOINTS")
	if exists {
		core.AppLog.Printf("Overriding ETCD : %s\n", c)
		f.EtcdEndpoints = f.EtcdEndpoints[:0]
		parts := strings.Split(c, ",")
		f.EtcdEndpoints = append(f.EtcdEndpoints, parts...)
	}
	cx := core.EtcdAtomic{Endpoints: f.EtcdEndpoints}
	lockPrefix := fmt.Sprintf("%s/node", f.Prefix)
	cnf := Config{}
	err = cx.Execute(lockPrefix, func(ctx core.Ctx) error {
		v, err := ctx.Get(f.NodeName)
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(v), &cnf)
		if err != nil {
			return err
		}
		cnf.Used = true
		err = ctx.Put(f.NodeName, string(util.ToJson(cnf)))
		if err != nil {
			return err
		}
		core.AppLog.Printf("config selected from etcd cluster : %s\n", f.NodeName)
		return nil
	})
	if err != nil {
		core.AppLog.Printf("error from etcd cluster : %s\n", err.Error())
		return nil
	}
	f.NodeId = int64(cnf.Sequence)
	core.AppLog.Printf("Overiding node id with %d\n", cnf.Sequence)
	f.HttpEndpoint = cnf.HttpEndpoint
	core.AppLog.Printf("Overiding http with %s\n", cnf.HttpEndpoint)
	f.Pgs.DatabaseURL = cnf.SqlEndpoint
	core.AppLog.Printf("Overiding sql with %s\n", cnf.SqlEndpoint)
	core.AppLog.Printf("Overiding tcp with %s\n", cnf.TcpEndpoint)
	parts := strings.Split(f.Evp.TcpEndpoint, ":")
	f.Evp.TcpEndpoint = parts[0] + "://" + cnf.TcpEndpoint + ":" + parts[2]
	core.AppLog.Printf("Using tcp endpoint : %s\n", f.Evp.TcpEndpoint)
	return nil
}
