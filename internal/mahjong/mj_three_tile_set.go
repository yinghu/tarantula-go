package mahjong

type ThreeTileSet struct {
	TileSetObj
}

func (f *ThreeTileSet) Append(t Tile) TileSet {
	f.Name = t.Suit
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *ThreeTileSet) Fallback(h *Hand) {
	if f.Size() == 2 && h.NextTile().Suit != f.Name {
		h.Formed = append(h.Formed, f.Formed())
		h.Push(h.NewTileSet(THREE_SET))
		return
	}
	for f.Size() > 0 {
		seq := h.NewTileSet(SEQ_SET)
		t := f.Head()
		seq.Append(t)
		h.Push(seq)
	}
}
