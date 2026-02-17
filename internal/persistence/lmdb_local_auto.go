package persistence

import (
	"gameclustering.com/internal/core"
	"wellquite.org/golmdb"
)

type LMDBLocalAuto struct {
	Path      string
	MapSizeMb int64
	Readers   uint
	MaxDbs    uint
	BatchSzie uint
	Env       *golmdb.LMDBClient
}

func (db *LMDBLocalAuto) Open() {
	var flag golmdb.EnvironmentFlag
	env, err := golmdb.NewLMDB(core.AppLog, db.Path, 0666, db.Readers, db.MaxDbs, flag, db.BatchSzie)
	if err != nil {
		panic(err)
	}
	db.Env = env
}
