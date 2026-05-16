package main

import (
	"gameclustering.com/internal/core"
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
