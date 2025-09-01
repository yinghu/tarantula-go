package mahjong

type FourTileSet struct {
	TileSetObj
}

func (f *FourTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *FourTileSet) Next(h *Hand,p Tile) (TileSet,error) {
	return h.NewTileSet(THREE_SET),nil
}

