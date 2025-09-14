package mahjong

import (
	"fmt"
	"slices"
	"strings"
)

type Hand struct {
	Formed     []Meld
	Tiles      []Tile
	Categories []int
	Flowers    []Tile
	Dots       []Tile
	Bamboos    []Tile
	Characters []Tile
	Hornors    []Tile
	Stack      []TileSet
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
	h.Dots = make([]Tile, 0)
	h.Bamboos = make([]Tile, 0)
	h.Characters = make([]Tile, 0)
	h.Hornors = make([]Tile, 0)
	h.Categories = []int{0, 0, 0, 0}
	h.Stack = make([]TileSet, 0)
}

func (h *Hand) Drop(drop Tile) error {
	deleted := false
	switch drop.Suit {
	case DOTS:
		for i := range h.Dots {
			if h.Dots[i] == drop {
				h.Dots = slices.Delete(h.Dots, i, i)
				deleted = true
				break
			}
		}
		if !deleted {
			return fmt.Errorf("drop not existed %v", drop)
		}
		h.Categories[0]--
	case BAMBOO:
		for i := range h.Bamboos {
			if h.Bamboos[i] == drop {
				h.Bamboos = slices.Delete(h.Bamboos, i, 1)
				deleted = true
				break
			}
		}
		if !deleted {
			return fmt.Errorf("drop not existed %v", drop)
		}
		h.Categories[1]--
	case CHARACTER:
		for i := range h.Characters {
			if h.Characters[i] == drop {
				h.Characters = slices.Delete(h.Characters, i, i)
				deleted = true
				break
			}
		}
		if !deleted {
			return fmt.Errorf("drop not existed %v", drop)
		}
		h.Categories[2]--
	case HORNOR:
		for i := range h.Hornors {
			if h.Hornors[i] == drop {
				h.Hornors = slices.Delete(h.Hornors, i, i)
				deleted = true
				break
			}
		}
		if !deleted {
			return fmt.Errorf("drop not existed %v", drop)
		}
		h.Categories[3]--
	}
	return nil
}

func (h *Hand) Draw(deck *Deck) error {
	t, err := deck.Draw()
	if err != nil {
		return err
	}
	return h.draw(t)
}

func (h *Hand) draw(t Tile) error {
	switch t.Suit {
	case FLOWER:
		h.Flowers = append(h.Flowers, t)
	case DOTS:
		h.Dots = append(h.Dots, t)
		h.Categories[0]++
	case BAMBOO:
		h.Bamboos = append(h.Bamboos, t)
		h.Categories[1]++
	case CHARACTER:
		h.Characters = append(h.Characters, t)
		h.Categories[2]++
	case HORNOR:
		h.Hornors = append(h.Hornors, t)
		h.Categories[3]++
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

func (h *Hand) MJ() bool {
	slices.SortFunc(h.Categories, func(a, b int) int {
		return a - b
	})
	fmt.Printf("%v\n",h.Categories)
	
	var eyeCount int
	var formed int
	for _, v := range h.Formed {
		fmt.Printf("X Set size : %v\n", v.Tiles)
		if v.Eye() {
			eyeCount++
		}
		formed++
	}
	return eyeCount == 1 && formed == 5 || formed == 14
}

func (h *Hand) Mahjong() bool {
	slices.SortFunc(h.Tiles, cmp)
	fmt.Printf("%v\n", h.Tiles)
	h.Push(h.NewTileSet(PONG))
	err := h.evaluate()
	if err != nil {
		fmt.Printf("no match %s\n", err.Error())
		return false
	}
	var eyeCount int
	var formed int
	for _, v := range h.Formed {
		fmt.Printf("X Set size : %v\n", v.Tiles)
		if v.Eye() {
			eyeCount++
		}
		formed++
	}
	return eyeCount == 1 && formed == 5 || formed == 14
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
					formed := tset.Formed()
					h.Formed = append(h.Formed, formed)

					if h.StackSize() == 0 {
						h.Push(h.NewTileSet(PONG))
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
	for h.StackSize() > 0 {
		merge := h.Pop()
		if merge.Size() > 0 {
			h.Formed = append(h.Formed, merge.Formed())
		}
	}
	return nil
}

func (h *Hand) NewTileSet(id uint8) TileSet {
	var tset TileSet
	switch id {
	case KNOG:
		tset = NewFourTileSet()
	case PONG:
		tset = NewThreeTileSet()
	case CHOW:
		tset = NewSequenceTileSet()
	case EYE:
		tset = NewTwoTileSet()
	}
	return tset
}

func (h *Hand) NextTile() Tile {
	return h.Tiles[0]
}

func (h *Hand) NextTiles(limit int) []Tile {
	lst := make([]Tile, 0)
	for i := range limit {
		if i < h.TileSize() {
			lst = append(lst, h.Tiles[i])
		}
	}
	return lst
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
