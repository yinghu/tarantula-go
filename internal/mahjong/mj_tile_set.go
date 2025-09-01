package mahjong

type TileSet interface {
	Full() bool
	Fallback(h *Hand) TileSet
	Append(t Tile) TileSet
	Head() Tile
	Tail() Tile
	Allowed(t Tile) bool
	Eye() bool
	Sequence() int
	Size() int
	Formed() Meld
	Next(h *Hand) TileSet
}

type TileSetObj struct {
	TileSet  []Tile
	FullSize int
	Seq      int
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
	return f.TileSet[sz-1] == t
}

func (f *TileSetObj) Eye() bool {
	return false
}

func (f *TileSetObj) Sequence() int {
	return f.Seq
}

func (f *TileSetObj) Size() int {
	return len(f.TileSet)
}
func (f *TileSetObj) Head() Tile {
	t := f.TileSet[0]
	f.TileSet = f.TileSet[1:]
	return t
}
func (f *TileSetObj) Tail() Tile {
	t := f.TileSet[f.Size()-1]
	f.TileSet = f.TileSet[:f.Size()-1]
	return t
}

func NewFourTileSet(sn int) TileSet {
	tset := FourTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 4
	tset.Seq = sn
	return &tset
}

func NewThreeTileSet(sn int) TileSet {
	tset := ThreeTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 3
	tset.Seq = sn
	return &tset
}

func NewSequenceTileSet(sn int) TileSet {
	tset := SequenceTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 3
	tset.Seq = sn
	return &tset
}

func NewTwoTileSet(sn int) TileSet {
	tset := TwoTileSet{}
	tset.TileSet = make([]Tile, 0)
	tset.FullSize = 2
	tset.Seq = sn
	return &tset
}
