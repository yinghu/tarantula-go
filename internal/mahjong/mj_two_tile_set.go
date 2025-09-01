package mahjong

import "fmt"

type TwoTileSet struct {
	TileSetObj
}

func (f *TwoTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)

	//fmt.Printf("PENDING EYE : %v\n", f.TileSet)

	return f
}

func (f *TwoTileSet) Eye() bool {
	return true
}

func (f *TwoTileSet) Next(h *Hand,p Tile) (TileSet, error) {
	return f, fmt.Errorf("no more match on two")
}
