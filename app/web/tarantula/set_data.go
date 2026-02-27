package main

import (
	"gameclustering.com/internal/core"
	badger "github.com/dgraph-io/badger/v4"
)

type SetData struct {
	Key     []byte
	Value   []byte
	Closing bool
	Msg     chan string
}

func (m *DataServiceProvider) runSetData() {
	start:
	for sd := range m.DSet {
		if sd.Closing {
			break
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
	core.AppLog.Debug().Msg("closing set operation")
	goto start
	
}
