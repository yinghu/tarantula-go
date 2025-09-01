package mahjong

import "fmt"

type TwoTileSet struct {
	TileSetObj
}

func (f *TwoTileSet) Fallback(h *Hand) TileSet {
	for f.Size() > 0 {
		h.PushTile(f.Tail())
	}
	return f
}
func (f *TwoTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	if len(f.TileSet) == 2 {
		fmt.Printf("EYE : %v\n", f.TileSet)
	}
	return f
}

func (f *TwoTileSet) Eye() bool {
	return true
}

func (f *TwoTileSet) Next(h *Hand) TileSet {
	return nil
}
