package mahjong

type ThreeTileSet struct {
	TileSetObj
}

func (f *ThreeTileSet) Fallback(tiles []Tile) TileSet {
	return f
}
func (f *ThreeTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}
