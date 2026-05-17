package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
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
	return nil
}

func (v *RepositoryObejctUpdate) confirm(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check confirm %v", t.Meta)
	return nil
}

func (v *RepositoryObejctUpdate) cancel(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check cancel %v", t.Meta)
	return nil
}
