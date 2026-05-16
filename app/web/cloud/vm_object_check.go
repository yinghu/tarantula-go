package main

import (
	"gameclustering.com/internal/core"
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
	return nil
}

func (v *VMObjectCheck) confirm(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check confirm %v", t.Meta)
	return nil
}

func (v *VMObjectCheck) cancel(t *protocol.Transaction) error {
	core.AppLog.Debug().Msgf("check cancel %v", t.Meta)
	return nil
}
