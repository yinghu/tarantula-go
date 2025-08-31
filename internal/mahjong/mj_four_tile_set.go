package mahjong

type FourTileSet struct {
	TileSetObj
}


func (f *FourTileSet) Fallback(h *Hand) TileSet {
	//tset := NewThreeTileSet(h.Sn)
	//h.Sn++
	
	return f
}
func (f *FourTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}
