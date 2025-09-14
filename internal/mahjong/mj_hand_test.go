package mahjong

import (
	"fmt"
	"slices"
	"testing"
)

func TestHand1(t *testing.T) {

	h := Hand{}
	h.New()

	b1 := Tile{Suit: "B", Rank: 1}
	b2 := Tile{Suit: "B", Rank: 2}
	b3 := Tile{Suit: "B", Rank: 3}

	b4 := Tile{Suit: "B", Rank: 2}
	b5 := Tile{Suit: "B", Rank: 3}
	b6 := Tile{Suit: "B", Rank: 4}

	b7 := Tile{Suit: "B", Rank: 3}
	b8 := Tile{Suit: "B", Rank: 4}
	b9 := Tile{Suit: "B", Rank: 5}

	c7 := Tile{Suit: "B", Rank: 4}
	c8 := Tile{Suit: "B", Rank: 5}
	c9 := Tile{Suit: "B", Rank: 6}

	p1 := Tile{Suit: "B", Rank: 1}
	p2 := Tile{Suit: "B", Rank: 1}

	tiles := []Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	for i := range tiles {
		h.draw(tiles[i])
	}
	if len(h.Bamboos) != 14 {
		t.Errorf("hand should be 14 %d", len(h.Bamboos))
	}
	if len(h.Characters) != 0 {
		t.Errorf("characters should be 0 %d", len(h.Characters))
	}
	if len(h.Dots) != 0 {
		t.Errorf("dots should be 0 %d", len(h.Dots))
	}
	if len(h.Hornors) != 0 {
		t.Errorf("hornors should be 0 %d", len(h.Hornors))
	}
	slices.SortFunc(h.Bamboos, cmp)
	fmt.Printf("%v\n", h.Bamboos)
	f4 := make([]TileSet, 0)
	f4 = append(f4, h.NewTileSet(KNOG))
	f3 := make([]TileSet, 0)
	f3 = append(f3, h.NewTileSet(PONG))
	f2 := make([]TileSet, 0)
	f2 = append(f2, h.NewTileSet(CHOW))
	tem := make([]TileSet, 0)
	for len(h.Bamboos) > 0 {
		t := h.Bamboos[0]
		h.Bamboos = h.Bamboos[1:]
		for len(f4) > 0 {
			ts := f4[0]
			f4 = f4[1:]
			if ts.Allowed(t) {
				ts.Append(t)
				if ts.Full() {
					c := ts.Formed()
					fmt.Printf("Kong %v\n", c)
					h.Formed = append(h.Formed, c)
				} else {
					tem = append(tem, ts)
				}
			} else {
				f := h.NewTileSet(KNOG)
				f.Append(t)
				tem = append(tem, f)
			}
		}
		f4 = append(f4, tem...)
		tem = tem[:0]
		for len(f3) > 0 {
			ts := f3[0]
			f3 = f3[1:]
			if ts.Allowed(t) {
				ts.Append(t)
				if ts.Full() {
					c := ts.Formed()
					fmt.Printf("Pong %v\n", c)
					h.Formed = append(h.Formed, c)
					tem = append(tem, ts)
				} else {
					tem = append(tem, ts)
				}
			} else {
				f := h.NewTileSet(PONG)
				f.Append(t)
				tem = append(tem, f)
			}
		}
		f3 = append(f3, tem...)
		tem = tem[:0]
		for len(f2) > 0 {
			ts := f2[0]
			f2 = f2[1:]
			if ts.Allowed(t) {
				ts.Append(t)
				if ts.Full() {
					c := ts.Formed()
					fmt.Printf("Chow %v\n", c)
					h.Formed = append(h.Formed, c)
					tem = append(tem, ts)
				}else{
					tem = append(tem, ts)
				}
			} else {
				tem = append(tem, ts)
				f := h.NewTileSet(CHOW)
				f.Append(t)
				tem = append(tem, f)
			}
		}
		f2 = append(f2, tem...)
		tem = tem[:0]
	}
	//for _, v := range h.Formed {
		//fmt.Printf("T Set size : %v\n", v)
	//}
}
