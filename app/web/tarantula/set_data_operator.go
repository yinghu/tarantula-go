package main

import (
	"gameclustering.com/internal/core"
	badger "github.com/dgraph-io/badger/v4"
)

const (
	SET_OPT_RECOVER int = 1
	SET_OPT_CLOSE   int = 2
	SET_OPT_CREATE  int = 3
	SET_OPT_UPDATE  int = 4
	SET_OPT_DELETE  int = 5

	SET_OPERATOR_NUM int = 8
)

type SetRes struct {
	Suc  bool
	Code int
	Err  error
}

func (m *DataServiceProvider) runSetData(num int) {
start:
	core.AppLog.Debug().Msgf("starting set operator %d", num)
	m.DWait.Wait()
	for sd := range m.DSet {
		if sd.Opt == SET_OPT_RECOVER {
			break
		} else if sd.Opt == SET_OPT_CLOSE {
			core.AppLog.Debug().Msgf("closing set data operator %d", num)
			return
		}
		err := m.Local.Db.Update(func(txn *badger.Txn) error {
			return txn.Set(sd.Data.Key, sd.Data.Value)
		})
		if err != nil {
			sd.Msg <- SetRes{Err: err}

		} else {
			sd.Msg <- SetRes{Suc: true}
		}
	}
	core.AppLog.Debug().Msgf("running recovery on operator %d", num)
	sync := <-m.DPull
	core.AppLog.Warn().Msgf("pulling data from %s", sync.Remote)
	for _, h := range sync.Hashs {
		core.AppLog.Warn().Msgf("recovering data from hash %d", h)
		//m.ClientPull()
	}
	m.DWait.Done()
	goto start
}
