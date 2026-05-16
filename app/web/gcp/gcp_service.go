package main

import (
	"os"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"
)

type GcpService struct {
	bootstrap.AppManager
	assetDir string
}

func (s *GcpService) Config() string {
	return "./gcp-conf.json"
}

func (s *GcpService) Start(f core.Env) error {
	s.AppManager.Start(f)
	s.Cluster().Subscribe(event.MESSAGE_TOPIC_NAME, &protocol.TopicEventListener{C: func() proto.Message {
		return &protocol.MessageEvent{}
	}, M: func(m proto.Message) {
		ro, ok := m.(*protocol.MessageEvent)
		if ok {
			core.AppLog.Debug().Msgf("MESSAGE event %s %s", ro.Message, ro.Source)
		} else {
			core.AppLog.Debug().Msg("wrong type")
		}
	}})
	s.Cluster().Register("update", &protocol.TccTransationListener{Reserve: func(e *protocol.Transaction) error {
		key, err := s.Cluster().AuthKey("gcp")
		if err != nil {
			return err
		}
		gcp := util.GcpComputeEngine{ServiceAccount: key.Gcp.Iam, ProjectId: "prismatic-grail-206205", Zone: "us-east1-c"}
		err = gcp.Auth()
		if err != nil {
			core.AppLog.Debug().Msgf("gcp auth error %s", err)
			return err
		}

		ins, err := gcp.Get("tarantula-build-02")
		if err != nil {
			core.AppLog.Debug().Msgf("gcp read error %s", err.Error())
			return err
		}

		ssh := util.SshClient{Host: ins.GetNetworkInterfaces()[0].AccessConfigs[0].GetNatIP(), User: "yinghu_lu", PrivateKey: key.Gcp.Ssh, KHFile: "../.ssh/known_hosts"}
		err = ssh.WithKey()
		if err != nil {
			core.AppLog.Debug().Msgf("gcp ssh error %s", err)
			return err
		}
		//var w bytes.Buffer

		gkey, err := s.Cluster().AuthKey("git")
		if err != nil {
			return err
		}
		err = os.WriteFile("id_ed25519", []byte(gkey.Git.Key), 0700)
		if err != nil {
			return err
		}
		f, err := os.Open("./id_ed25519")
		//err = ssh.Run("ssh-keyscan -t ed25519 github.com >> .ssh/known_hosts", &w)
		if err != nil {
			return err
		}
		err = ssh.Upload(f, "/home/yinghu_lu/.ssh/id_ed25519", "0700") //perm 0600
		if err != nil {
			core.AppLog.Debug().Msgf("scp ssh error %s", err)
			return err
		}
		//core.AppLog.Debug().Msgf("git known host added :%s", w.String())
		f.Close()
		gcp.Close()
		ssh.Close()
		return nil
	}, Confirm: func(e *protocol.Transaction) error {
		return nil
	}, Cancel: func(e *protocol.Transaction) error {
		return nil
	}})
	core.AppLog.Printf("Asset service started %s %s\n", f.HttpBinding, s.assetDir)
	return nil
}
