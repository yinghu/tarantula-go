package main

import (
	"fmt"
	"net/http"

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

		//cmd := exec.Command("ssh-keyscan", "-H", ins.GetNetworkInterfaces()[0].AccessConfigs[0].GetNatIP())
		//out, err := cmd.Output()
		//if err != nil {
		//core.AppLog.Debug().Msgf("ssh-kenscann error %s", err.Error())
		//return err
		//}
		//err = os.WriteFile("../.ssh/known_host", out, 0644)
		//if err != nil {
		//core.AppLog.Debug().Msgf("write file error %s", err.Error())
		//return err
		//}
		ssh := util.SshClient{Host: ins.GetNetworkInterfaces()[0].AccessConfigs[0].GetNatIP(), User: "yinghu_lu", PrivateKey: key.Gcp.Ssh, KHFile: "../.ssh/known_hosts"}
		err = ssh.WithKey()
		if err != nil {
			core.AppLog.Debug().Msgf("gcp ssh error %s", err)
			return err
		}
		core.AppLog.Debug().Msgf("gcp ssh ok %s", err)
		ssh.Run("pwd")
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
