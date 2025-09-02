package mahjong


type ThreeTileSet struct {
	TileSetObj
}

func (f *ThreeTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *ThreeTileSet) Fallback(h *Hand) {
	for f.Size() > 0 {
		seq := h.NewTileSet(SEQ_SET)
		t := f.Head()
		seq.Append(t)
		h.Push(seq)
	}
}
