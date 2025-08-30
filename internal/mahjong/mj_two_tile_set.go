package mahjong

type TwoTileSet struct {
	TileSetObj
}

func (f *TwoTileSet) Fallback(tiles []Tile) TileSet {
	return f
}
func (f *TwoTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *TwoTileSet) Eye() bool {
	return true
}
