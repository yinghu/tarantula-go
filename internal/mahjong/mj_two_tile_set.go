package mahjong


type TwoTileSet struct {
	TileSetObj
}

func (f *TwoTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	return f
}

func (f *TwoTileSet) Eye() bool {
	return true
}


func (f *TwoTileSet) Fallback(h *Hand) {

}
