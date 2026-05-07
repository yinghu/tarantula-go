package main

import (
	"fmt"
	"net/http"
	"os"

	"gameclustering.com/internal/bootstrap"
	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/protocol"
	"google.golang.org/protobuf/proto"
)

type AdminService struct {
	bootstrap.AppManager
	assetDir    string
	contentDir  string
	publishDir  string
	managedApps []string
}

func (s *AdminService) Config() string {
	return "/etc/tarantula/admin-conf.json"
}

func (s *AdminService) Start(f core.Env) error {
	f.AuthLevel = core.ADMIN_ACCESS_CONTROL
	s.AppManager.Start(f)
	s.managedApps = f.ManagedApps
	s.contentDir = fmt.Sprintf("%s/%s", f.HomeDir, "bin")
	s.assetDir = fmt.Sprintf("%s/%s/%s", f.HomeDir, f.GroupName, "asset")
	os.MkdirAll(s.assetDir, 0755)
	s.publishDir = s.contentDir + "/tarantula"
	err := s.createSchema()
	if err != nil {
		return err
	}
	hash, err := s.Authenticator().HashPassword("password")
	if err != nil {
		return err
	}
	err = s.SaveLogin(&protocol.LoginObject{Name: "root", Password: hash, AccessControl: uint32(core.SUDO_ACCESS_CONTROL)})
	if err != nil {
		core.AppLog.Debug().Msg("Root already existed")
	}
	s.Cluster().Subscribe(event.TASK_TOPIC_NAME, &protocol.TopicEventListener{C: func() proto.Message {
		return &protocol.TaskEvent{}
	}, M: func(m proto.Message) {
		ro, ok := m.(*protocol.TaskEvent)
		if ok {
			core.AppLog.Debug().Msgf("validator event %v", ro.Validator)
			core.AppLog.Debug().Msgf("transaction event %v", ro.Job)
		} else {
			core.AppLog.Debug().Msg("wrong type")
		}
	}})
	http.Handle("/admin/webprotected/{name}", bootstrap.Logging(&AdminWebProtected{AdminService: s}))
	http.Handle("/admin/web/{name}", bootstrap.Logging(&AdminWebIndex{AdminService: s}))
	//handle / context from nginx proxy
	http.Handle("/admin/{name}", bootstrap.Logging(&AdminWebIndex{AdminService: s}))

	http.Handle("/admin/cs/message/send", bootstrap.Logging(&CSMessager{AdminService: s}))

	http.Handle("/admin/cs/query/topic/{topic}", bootstrap.Logging(&CSQueryTopic{AdminService: s}))
	http.Handle("/admin/cs/query/object/{topic}", bootstrap.Logging(&CSQueryObject{AdminService: s}))
	http.Handle("/admin/cs/inventory/grant", bootstrap.Logging(&CSGranter{AdminService: s}))
	http.Handle("/admin/cs/inventory/load", bootstrap.Logging(&CSInventoryLoader{AdminService: s}))


	http.Handle("/admin/view/{id}", bootstrap.Logging(&AdminItemViewer{AdminService: s}))
	http.Handle("/admin/repo/sync", bootstrap.Logging(&AdminPublisher{AdminService: s}))
	http.Handle("/admin/env", bootstrap.Logging(&AdminEnv{AdminService: s}))
	http.Handle("/admin/snowflake/parse", bootstrap.Logging(&AdminParseSnowFlakeId{AdminService: s}))
	http.Handle("/admin/login/add", bootstrap.Logging(&SudoAddLogin{AdminService: s}))
	http.Handle("/admin/password", bootstrap.Logging(&AdminChangePwd{AdminService: s}))
	http.Handle("/admin/accesskey", bootstrap.Logging(&AdminCreateAccessKey{AdminService: s}))
	http.Handle("/admin/login", bootstrap.Logging(&AdminLogin{AdminService: s}))

	http.Handle("/admin/presence/hashring", bootstrap.Logging(&AdminHashRingEndpoint{AdminService: s}))
	http.Handle("/admin/presence/keyring/{key}", bootstrap.Logging(&AdminKeyRingEndpoint{AdminService: s}))
	http.Handle("/admin/presence/subscription/task", bootstrap.Logging(&AdminSubscriptionTaskEndpoint{AdminService: s}))
	http.Handle("/admin/presence/subscription/topic", bootstrap.Logging(&AdminSubscriptionTopicEndpoint{AdminService: s}))

	http.Handle("/admin/cluster/delete/{key}", bootstrap.Logging(&AdminClusterDelete{AdminService: s}))
	http.Handle("/admin/cluster/reset", bootstrap.Logging(&AdminClusterReset{AdminService: s}))

	core.AppLog.Info().Msgf("Admin service started %s\n", f.HttpBinding)
	return nil
}
