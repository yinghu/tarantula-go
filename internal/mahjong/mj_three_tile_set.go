package mahjong

type ThreeTileSet struct {
	TileSetObj
}

func (f *ThreeTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *ThreeTileSet) Next(h *Hand,p Tile) (TileSet, error) {
	if f.Size() == 1 {
		seq := h.NewTileSet(SEQ_SET)
		seq.Append(f.Head())
		return seq, nil
	}
	two := h.NewTileSet(TWO_SET)
	for f.Size() > 0{
		two.Append(f.Head())
	}
	return two, nil
}
