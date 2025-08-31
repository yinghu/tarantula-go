package mahjong

type SequenceTileSet struct {
	TileSetObj
}

func (f *SequenceTileSet) Fallback(h *Hand) TileSet {
	tset := h.NewTileSet(TWO_SET)
	for f.Size() > 0 {
		h.PushTile(f.Pop())
	}
	return tset.Append(h.PopTile())
}
func (f *SequenceTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *SequenceTileSet) Allowed(t Tile) bool {
	sz := len(f.TileSet)
	if sz == 0 {
		return true
	}
	if sz == f.FullSize {
		return false
	}
	if f.TileSet[sz-1].Suit == t.Suit && f.TileSet[sz-1].Rank+1 == t.Rank {
		return true
	}
	return false
}
