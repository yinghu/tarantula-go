package mahjong

type TwoTileSet struct {
	TileSetObj
}

func (f *TwoTileSet) Append(t Tile) TileSet {
	f.Name = t.Suit
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *TwoTileSet) Eye() bool {
	return true
}

func (f *TwoTileSet) Fallback(h *Hand) {

}
