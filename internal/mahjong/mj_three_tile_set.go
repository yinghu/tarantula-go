package mahjong

type ThreeTileSet struct {
	TileSetObj
}

func (f *ThreeTileSet) Fallback(h *Hand) TileSet {
	tset := h.NewTileSet(SEQ_SET)
	for f.Size() > 0 {
		h.PushTile(f.Head())
	}
	return tset.Append(h.PopTile())
}
func (f *ThreeTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *ThreeTileSet) Next(h *Hand) TileSet {
	return h.NewTileSet(SEQ_SET)
}
