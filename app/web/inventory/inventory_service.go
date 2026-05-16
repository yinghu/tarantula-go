package main

import (
	"fmt"
	"net/http"
	"os"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
	"google.golang.org/protobuf/proto"
)

type InventoryService struct {
	bootstrap.AppManager
}

func (s *InventoryService) Config() string {
	return "./inventory-conf.json"
}

func (s *InventoryService) Start(f core.Env) error {
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
	s.Cluster().Register("grant", &protocol.TccTransationListener{Reserve: func(e *protocol.Transaction) error {
		mf := persistence.NewCommodityObjectFactory()
		core.AppLog.Debug().Msgf("SYSTEMID %d", s.ToSystemId(e.Object.Key.Array))
		core.AppLog.Debug().Msgf("META %v", e.Meta)
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
		f, err := os.OpenFile("id_ed25519", os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		f.Write([]byte(gkey.Git.Key))
		//err = ssh.Run("ssh-keyscan -t ed25519 github.com >> .ssh/known_hosts", &w)
		//if err != nil {
		//return err
		//}
		err = ssh.Upload(f, "home/yinghu_lu/.ssh/id_ed25519", "0600") //perm 0600
		if err != nil {
			core.AppLog.Debug().Msgf("scp ssh error %s", err)
			return err
		}
		//core.AppLog.Debug().Msgf("git known host added :%s", w.String())
		f.Close()
		gcp.Close()
		ssh.Close()
		obj, err := mf.Message(e.Object)
		if err != nil {
			core.AppLog.Warn().Msgf("wrong decode format %s", err)
			return err
		}
		m, ok := obj.(*protocol.Commodity)
		if !ok {
			return fmt.Errorf("wrong data passed from task")
		}
		core.AppLog.Debug().Msgf("%v", m)
		return nil
	}, Confirm: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("C META %v", e.Meta)
		return nil
	}, Cancel: func(e *protocol.Transaction) error {
		core.AppLog.Debug().Msgf("N META %v", e.Meta)
		return nil
	}})
	http.Handle("/inventory/cluster/update", bootstrap.Logging(&InventoryClusterUpdate{InventoryService: s}))
	return nil
}
