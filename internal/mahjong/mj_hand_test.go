package mahjong

import "testing"

func TestHandTile(t *testing.T) {
	h := Hand{}
	h.New()
	t1 := Tile{Suit: "D", Rank: 1}
	t2 := Tile{Suit: "D", Rank: 2}
	t3 := Tile{Suit: "D", Rank: 3}
	t4 := Tile{Suit: "D", Rank: 4}
	t5 := Tile{Suit: "D", Rank: 5}
	h.PushTile(t1)
	h.PushTile(t2)
	h.PushTile(t3)
	h.PushTile(t4)
	h.PushTile(t5)
	if h.TileSize() != 5 {
		t.Errorf("tile size should be 5 %d", h.TileSize())
	}
	if h.Tiles[0] != t5 {
		t.Errorf("first should be t5 %v", h.Tiles[0])
	}
	if h.Tiles[4] != t1 {
		t.Errorf("last should be t1 %v", h.Tiles[4])
	}
	p := h.PopTile()
	if h.TileSize() != 4 {
		t.Errorf("tile size should be 4 %d", h.TileSize())
	}
	if p != t5 {
		t.Errorf("pop should be t5 %v", p)
	}
	h.PopTile()
	h.PopTile()
	h.PopTile()
	h.PopTile()
	if h.TileSize() != 0 {
		t.Errorf("tile size should be 0 %d", h.TileSize())
	}
	h.PushTile(t1)
	h.PushTile(t2)
	h.PushTile(t3)
	h.PushTile(t4)
	h.PushTile(t5)
	if h.TileSize() != 5 {
		t.Errorf("tile size should be 5 %d", h.TileSize())
	}
}

func TestHandTileSet(t *testing.T) {
	h := Hand{}
	h.New()
	t1 := h.NewTileSet(FOUR_SET)
	t2 := h.NewTileSet(THREE_SET)
	t3 := h.NewTileSet(SEQ_SET)
	t4 := h.NewTileSet(TWO_SET)
	if t1.Sequence() != 1 {
		t.Errorf("Seq should be 1 %d", t1.Sequence())
	}
	if t2.Sequence() != 2 {
		t.Errorf("Seq should be 2 %d", t2.Sequence())
	}
	if t3.Sequence() != 3 {
		t.Errorf("Seq should be 3 %d", t3.Sequence())
	}
	if t4.Sequence() != 4 {
		t.Errorf("Seq should be 4 %d", t4.Sequence())
	}
	h.PushTileSet(t1)
	h.PushTileSet(t2)
	h.PushTileSet(t3)
	h.PushTileSet(t4)
	if h.TileSetSize() != 4 {
		t.Errorf("tile set size should be 4 %d", h.TileSetSize())
	}
	t0 := h.Pending[0]
	p := h.PopTileSet()
	if t0.Sequence() != p.Sequence() && t0.Sequence() != 4 {
		t.Errorf("first one should be t4 %d %d", t0.Sequence(), p.Sequence())
	}
	if h.TileSetSize() != 3 {
		t.Errorf("tile set size should be 3 %d", h.TileSetSize())
	}
	h.PopTileSet()
	h.PopTileSet()
	h.PopTileSet()
	if h.TileSetSize() != 0 {
		t.Errorf("tile set size should be 0 %d", h.TileSetSize())
	}
}

func TestHandFourTileSet(t *testing.T) {
	h := Hand{}
	h.New()
	t1 := Tile{Suit: "D", Rank: 1}
	t2 := Tile{Suit: "D", Rank: 1}
	t3 := Tile{Suit: "D", Rank: 1}
	t4 := Tile{Suit: "D", Rank: 1}

	t4set := h.NewTileSet(FOUR_SET)
	t4set.Append(t1)
	t4set.Append(t2)
	t4set.Append(t3)
	t4set.Append(t4)
	if t4set.Size()!=4{
		t.Errorf("size should be 4 %d",t4set.Size())
	}
	t3set := t4set.Fallback(&h)
	if t3set.Size() !=1 {
		t.Errorf("size should be 1 %d",t3set.Size())
	}
	if t4set.Size()!=0{
		t.Errorf("size should be 0 %d",t4set.Size())
	}
	if h.TileSize()!=3{
		t.Errorf("size should be 3 %d",h.TileSize())
	}
}
