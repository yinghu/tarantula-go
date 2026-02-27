package main

import (
	"gameclustering.com/internal/core"
	badger "github.com/dgraph-io/badger/v4"
)

const (
	SET_OPT_RECOVER int = 1
	SET_OPT_CLOSE   int = 2

	SET_OPERATOR_NUM int = 3
)

type SetRes struct {
	Suc  bool
	Code int
	Err  error
}

type SetData struct {
	Key   []byte
	Value []byte
	Opt   int
	Msg   chan SetRes
}

func (m *DataServiceProvider) runSetData(num int) {
start:
	core.AppLog.Debug().Msgf("starting set operation on %d",num)
	m.DWait.Wait()
	for sd := range m.DSet {
		if sd.Opt == SET_OPT_RECOVER {
			break
		} else if sd.Opt == SET_OPT_CLOSE {
			core.AppLog.Debug().Msg("closing set data operation")
			return
		}
		err := m.Local.Db.Update(func(txn *badger.Txn) error {
			return txn.Set(sd.Key, sd.Value)
		})
		if err != nil {
			sd.Msg <- SetRes{Err: err}

		} else {
			sd.Msg <- SetRes{Suc: true}
		}
	}
	core.AppLog.Debug().Msg("running recover operation")
	sync := <-m.DPull
	core.AppLog.Warn().Msgf("pulling data from %s", sync.Remote)
	for _, h := range sync.Hashs {
		core.AppLog.Warn().Msgf("recovering data from hash %d", h)
	}
	m.DWait.Done()
	goto start

}
