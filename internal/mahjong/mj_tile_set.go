package mahjong

type TileSet interface {
	Full() bool
	Append(t Tile) TileSet
	Head() Tile
	Allowed(t Tile) bool
	Size() int
	Formed() Meld
	Fallback(h *Hand)
	Suit() string
}

type TileSetObj struct {
	TileSet  []Tile
	FullSize int
	Name     string
}

func (f *TileSetObj) Full() bool {
	return len(f.TileSet) == f.FullSize
}

func (f *TileSetObj) Formed() Meld {
	m := Meld{Tiles: f.TileSet}
	f.TileSet = f.TileSet[:0]
	return m
}

func (f *TileSetObj) Allowed(t Tile) bool {
	sz := len(f.TileSet)
	if sz == 0 {
		return true
	}
	if sz == f.FullSize {
		return false
	}
	return f.TileSet[sz-1] == t
}

func (f *TileSetObj) Size() int {
	return len(f.TileSet)
}
func (f *TileSetObj) Head() Tile {
	t := f.TileSet[0]
	f.TileSet = f.TileSet[1:]
	return t
}

func (f *TileSetObj) Suit() string {
	return f.Name
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
