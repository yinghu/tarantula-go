package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
	"gameclustering.com/internal/util"
)

func NewRepositoryObejctUpdate(s *CloudService) *protocol.TccTransationListener {
	vm := RepositoryObejctUpdate{s}
	tcc := protocol.TccTransationListener{}
	tcc.Reserve = vm.reserse
	tcc.Confirm = vm.confirm
	tcc.Cancel = vm.cancel
	return &tcc
}

type RepositoryObejctUpdate struct {
	*CloudService
}

func (v *RepositoryObejctUpdate) reserse(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("update reserve %v", t.Meta)
	tf := persistence.NewRepositoryObjectFactory()
	obj, err := tf.Message(t.Object)
	if err != nil {
		return err
	}
	vm, ok := obj.(*protocol.RepositoryObject)
	if !ok {
		return fmt.Errorf("wrong message type")
	}
	core.AppLog.Debug().Msgf("repository object %v", vm)
	github, err := v.AppManager.Cluster().AuthKey("git")
	if err != nil {
		core.AppLog.Debug().Msgf("no git key %s", err.Error())
		return v.insert(t.Meta)
	}
	gapi := util.GitHubApi{Token: github.Git.Token, Org: github.Git.Org}
	repos, err := gapi.ListRepos()
	if err != nil {
		core.AppLog.Debug().Msgf("repo error %s", err.Error())
	}
	core.AppLog.Debug().Msgf("repos %v", repos)
	return v.insert(t.Meta)
}

func (v *RepositoryObejctUpdate) confirm(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check confirm %v", t.Meta)
	return v.insert(t.Meta)
}

func (v *RepositoryObejctUpdate) cancel(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check cancel %v", t.Meta)
	return v.insert(t.Meta)
}
