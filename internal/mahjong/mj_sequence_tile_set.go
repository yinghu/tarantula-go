package mahjong

import "fmt"

type SequenceTileSet struct {
	TileSetObj
}

func (f *SequenceTileSet) Append(t Tile) TileSet {
	f.TileSet = append(f.TileSet, t)
	//fmt.Printf("PENDING CHOW : %v\n", f.TileSet)
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

func (f *SequenceTileSet) Next(h *Hand, p Tile) (TileSet, error) {
	if f.Size() == 1 {
		two := h.NewTileSet(TWO_SET)
		two.Append(f.Head())
		return two, nil
	}
	tail := f.Tail()
	if tail == p {
		h.PushTile(f.Head())
		two := h.NewTileSet(TWO_SET)
		two.Append(tail)
		return two, nil
	}
	return f, fmt.Errorf("no more match")
}
