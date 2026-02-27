package main

import (
	"gameclustering.com/internal/core"
	badger "github.com/dgraph-io/badger/v4"
)

const (
	SET_OPT_RECOVER int = 1
	SET_OPT_CLOSE   int = 2
)

type SetData struct {
	Key   []byte
	Value []byte
	Opt   int
	Msg   chan string
}

func (m *DataServiceProvider) runSetData() {
start:
	core.AppLog.Debug().Msg("starting set operation")
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
			sd.Msg <- err.Error()

		} else {
			sd.Msg <- "saved"
		}
	}
	core.AppLog.Debug().Msg("running recover operation")
	sync := <-m.DPull
	core.AppLog.Warn().Msgf("pulling data from %s", sync.Remote)
	for _, h := range sync.Hashs {
		core.AppLog.Warn().Msgf("recovering data from hash %d", h)
	}
	goto start

}
