package mahjong

type TileSet interface {
	Full() bool
	Fallback(tiles []Tile) TileSet
	Append(t Tile) TileSet
	Allowed(t Tile) bool
	Eye() bool
	Debug() []Tile
}

type TileSetObj struct {
	TileSet  []Tile
	FullSize int
}

func (f *TileSetObj) Full() bool {
	return len(f.TileSet) == f.FullSize
}

func (f *TileSetObj) Debug() []Tile {
	return f.TileSet
}

func (f *TileSetObj) Allowed(t Tile) bool {
	sz := len(f.TileSet)
	if sz == 0 {
		return true
	}
	return f.TileSet[sz-1] == t
}

func (f *TileSetObj) Eye() bool {
	return false
}

func NewFourTileSet() TileSet {
	tset := FourTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 4
	return &tset
}

func NewThreeTileSet() TileSet {
	tset := ThreeTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 3
	return &tset
}

func NewSequenceTileSet() TileSet {
	tset := SequenceTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 3
	return &tset
}

func NewTwoTileSet() TileSet {
	tset := TwoTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 2
	return &tset
}
