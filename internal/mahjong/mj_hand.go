package mahjong

import (
	"fmt"
	"slices"
	"strings"
)

const (
	FOUR_SET  int = 0
	THREE_SET int = 1
	SEQ_SET   int = 2
	TWO_SET   int = 3
)

type Hand struct {
	Formed  []Meld
	Tiles   []Tile
	Flowers []Tile

	Stack []TileSet
}

func cmp(a, b Tile) int {
	if a.Suit == b.Suit {
		diff := a.Rank - b.Rank
		return int(diff)
	}
	return strings.Compare(a.Suit, b.Suit)
}

func (h *Hand) New() {
	h.Formed = make([]Meld, 0)
	h.Tiles = make([]Tile, 0)
	h.Flowers = make([]Tile, 0)
	h.Stack = make([]TileSet, 0)
}

func (h *Hand) Draw(deck *Deck) error {
	t, err := deck.Draw()
	if err != nil {
		return err
	}
	switch t.Suit {
	case FLOWER:
		h.Flowers = append(h.Flowers, t)
	default:
		h.Tiles = append(h.Tiles, t)
		slices.SortFunc(h.Tiles, cmp)
	}
	return nil
}

func (h *Hand) Knog(deck *Deck) error {
	t, err := deck.Kong()
	if err != nil {
		return err
	}
	switch t.Suit {
	case FLOWER:
		h.Flowers = append(h.Flowers, t)
	default:
		h.Tiles = append(h.Tiles, t)
		slices.SortFunc(h.Tiles, cmp)
	}
	return nil
}

func (h *Hand) Mahjong() bool {
	slices.SortFunc(h.Tiles, cmp)
	fmt.Printf("%v\n", h.Tiles)
	h.Push(h.NewTileSet(THREE_SET))
	err := h.evaluate()
	if err != nil {
		fmt.Printf("no match %s\n", err.Error())
		return false
	}
	h.merge()
	var eyeCount int
	var formed int
	for _, v := range h.Formed {
		fmt.Printf("X Set size : %v\n", v.Tiles)
		if v.Eye() {
			eyeCount++
		}
		formed++
	}
	return eyeCount == 1 && formed == 5
}
func (h *Hand) evaluate() error {
	tem := make([]TileSet, 0)
	for h.TileSize() > 0 {
		t := h.PopTile()
		for h.StackSize() > 0 {
			tset := h.Pop()
			if tset.Allowed(t) {
				tset.Append(t)
				if tset.Full() {
					h.Formed = append(h.Formed, tset.Formed())
					if h.StackSize() == 0 {
						h.Push(h.NewTileSet(THREE_SET))
					}
				} else {
					h.Push(tset)
				}
				slices.Reverse(tem)
				for _, v := range tem {
					h.Push(v)
				}
				tem = tem[:0]
				break
			}
			if h.StackSize() > 0 {
				tem = append(tem, tset)
				continue
			}
			slices.Reverse(tem)
			for _, v := range tem {
				h.Push(v)
			}
			tem = tem[:0]
			tset.Fallback(h)
		}
	}
	return nil
}

func (h *Hand) merge() {
	if h.StackSize() == 1 {
		h.Formed = append(h.Formed, h.Pop().Formed())
		return
	}
	tiles := make([]Tile, 0)
	for h.StackSize() > 0 {
		meld := h.Pop().Formed()
		tiles = append(tiles, meld.Tiles...)
	}
	slices.SortFunc(tiles, cmp)
	fmt.Printf("Merge %v\n",tiles)
	//th := Hand{}
	//th.New()
	//th.Tiles = tiles
	//slices.SortFunc(th.Tiles, cmp)
	//h.Push(h.NewTileSet(THREE_SET))
	//th.evaluate()
	//h.Formed = append(h.Formed, th.Formed...)
}

func (h *Hand) NewTileSet(id int) TileSet {
	var tset TileSet
	switch id {
	case FOUR_SET:
		tset = NewFourTileSet()
	case THREE_SET:
		tset = NewThreeTileSet()
	case SEQ_SET:
		tset = NewSequenceTileSet()
	case TWO_SET:
		tset = NewTwoTileSet()
	}
	return tset
}

func (h *Hand) PopTile() Tile {
	t := h.Tiles[0]
	h.Tiles = h.Tiles[1:]
	return t
}

func (h *Hand) PushTile(t Tile) {
	h.Tiles = slices.Insert(h.Tiles, 0, t)
}

func (h *Hand) TileSize() int {
	return len(h.Tiles)
}

func (h *Hand) Push(ts TileSet) {
	h.Stack = append(h.Stack, ts)
}

func (h *Hand) Pop() TileSet {
	ts := h.Stack[h.StackSize()-1]
	h.Stack = h.Stack[:h.StackSize()-1]
	return ts
}

func (h *Hand) StackSize() int {
	return len(h.Stack)
}
