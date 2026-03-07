package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Sql struct {
	Enabled     bool   `json:"Enabled"`
	DatabaseURL string `json:"DatabaseURL"`
}

type EventEndpoint struct {
	Enabled         bool   `json:"Enabled"`
	OutboundEnabled bool   `json:"OutboundEnabled"`
	TcpEndpoint     string `json:"TcpEndpoint"`
}

type Env struct {
	Prefix        string        `json:"Prefix"`
	Standalone    bool          `json:"Standalone"`
	GroupName     string        `json:"GroupName"`
	NodeName      string        `json:"NodeName"`
	NodeId        int64         `json:"NodeId"`
	Host          string        `json:"Host"`
	HttpBinding   string        `json:"HttpBinding"`
	HttpEndpoint  string        `json:"HttpEndpoint"`
	Evp           EventEndpoint `json:"EventEndpoint"`
	EtcdEndpoints []string      `json:"EtcdEndpoints"`
	ManagedApps   []string      `json:"ManagedApps"`
	Pgs           Sql           `json:"Sql"`
	HomeDir       string        `json:"HomeDir"`
	LogTruncated  bool          `json:"LogTruncated"`
	AuthLevel     int32         `json:"AuthLevel"`
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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	mountDir := fmt.Sprintf("%s/%s", homeDir, f.GroupName)
	f.HomeDir = homeDir
	err = os.MkdirAll(mountDir, 0755)
	if err != nil {
		return err
	}
	CreateAppLog(mountDir, f.LogTruncated)
	f.Prefix = "dev"

	c, exists := os.LookupEnv("HOST")
	if exists {
		f.Host = c
		f.HttpEndpoint = c
		parts := strings.Split(f.Evp.TcpEndpoint, ":")
		f.Evp.TcpEndpoint = parts[0] + "://" + c + ":" + parts[2]
	}

	c, exists = os.LookupEnv("ENV")
	if exists {
		f.Prefix = c
	}

	c, exists = os.LookupEnv("SEQ")
	if exists {
		seq, err := strconv.Atoi(c)
		if err != nil {
			return err
		}
		f.NodeId = int64(f.NodeId + int64(seq))
		f.NodeName = fmt.Sprintf("%s.%d", f.GroupName, f.NodeId)

	}
	c, exists = os.LookupEnv("ETCD_ENDPOINTS")
	if exists {
		f.EtcdEndpoints = f.EtcdEndpoints[:0]
		parts := strings.Split(c, ",")
		f.EtcdEndpoints = append(f.EtcdEndpoints, parts...)
	}

	c, exists = os.LookupEnv("SQL_ENDPOINT")
	if exists {
		f.Pgs.DatabaseURL = c
	}
	AppLog.Debug().Msgf("CONF : %s %s %s %d %s %s %s %s", f.Prefix, f.GroupName, f.NodeName, f.NodeId, f.HttpEndpoint, f.Evp.TcpEndpoint, f.EtcdEndpoints[0], f.Pgs.DatabaseURL)

	return nil
}
