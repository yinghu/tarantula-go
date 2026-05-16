package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
)

func NewVMObejctCreate(s *CloudService) *protocol.TccTransationListener {
	vm := VMObjectCreate{s}
	tcc := protocol.TccTransationListener{}
	tcc.Reserve = vm.reserse
	tcc.Confirm = vm.confirm
	tcc.Cancel = vm.cancel
	return &tcc
}

type VMObjectCreate struct {
	*CloudService
}

func (v *VMObjectCreate) reserse(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("create reserve %v", t.Meta)
	tf := persistence.NewVMObjectFactory()
	obj, err := tf.Message(t.Object)
	if err != nil {
		return err
	}
	vm, ok := obj.(*protocol.VMObject)
	if !ok {
		return fmt.Errorf("wrong message type")
	}
	core.AppLog.Debug().Msgf("vm object %v", vm)
	return nil
}

func (v *VMObjectCreate) confirm(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("create confirm %v", t.Meta)
	return nil
}

func (v *VMObjectCreate) cancel(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("create cancel %v", t.Meta)

	return nil
}
