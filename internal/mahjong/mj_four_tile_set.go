package mahjong

type FourTileSet struct {
	TileSetObj
}


func (f *FourTileSet) Fallback(h *Hand) TileSet {
	tset := h.NewTileSet(THREE_SET)
	
	return tset
}
func (f *FourTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}
