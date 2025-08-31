package mahjong

type TwoTileSet struct {
	TileSetObj
}

func (f *TwoTileSet) Fallback(h *Hand) TileSet {
	for f.Size() > 0 {
		//h.PushTile(f.Pop())
	}
	return f
}
func (f *TwoTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *TwoTileSet) Eye() bool {
	return true
}
