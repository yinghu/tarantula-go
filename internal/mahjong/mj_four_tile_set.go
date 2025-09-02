package mahjong

type FourTileSet struct {
	TileSetObj
}

func (f *FourTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *FourTileSet) Fallback(h *Hand) {

}