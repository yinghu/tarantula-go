package mahjong

type SequenceTileSet struct {
	TileSetObj
}

func (f *SequenceTileSet) Append(t Tile) TileSet {
	f.Name = t.Suit
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *SequenceTileSet) Allowed(t Tile) bool {
	if t.Suit == HORNOR || t.Suit == FLOWER {
		return false
	}
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

func (f *SequenceTileSet) Fallback(h *Hand) {
	h.Push(h.NewTileSet(PONG))
	h.Push(f)
}
