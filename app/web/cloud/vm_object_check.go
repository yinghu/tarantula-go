package main

import (
	"fmt"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
)

func NewVMObejctCheck(s *CloudService) *protocol.TccTransationListener {
	vm := VMObjectCheck{s}
	tcc := protocol.TccTransationListener{}
	tcc.Reserve = vm.reserse
	tcc.Confirm = vm.confirm
	tcc.Cancel = vm.cancel
	return &tcc
}

type VMObjectCheck struct {
	*CloudService
}

func (v *VMObjectCheck) reserse(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check reserve %v", t.Meta)
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
	return v.insert(t.Meta)
}

func (v *VMObjectCheck) confirm(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check confirm %v", t.Meta)
	return v.insert(t.Meta)
}

func (v *VMObjectCheck) cancel(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check cancel %v", t.Meta)
	return v.insert(t.Meta)
}
