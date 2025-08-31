package mahjong

type FourTileSet struct {
	TileSetObj
}

func (f *FourTileSet) Fallback(h *Hand) TileSet {
	tset := h.NewTileSet(THREE_SET)
	for f.Size() > 0 {
		h.PushTile(f.Pop())
	}
	return tset.Append(h.PopTile())
}
func (f *FourTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}
