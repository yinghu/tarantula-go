package persistence

type LMDBLocalAuto struct{
	Path      string
	MapSizeMb int64
	Readers   int
}

func (db *LMDBLocalAuto) Open() {

}