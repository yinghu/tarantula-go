package mahjong

import (
	"fmt"
	"slices"
	"testing"
)

func TestTileSort(t *testing.T) {
	d1 := Tile{Suit: "D", Num: 9}
	d2 := Tile{Suit: "D", Num: 2}
	d5 := Tile{Suit: "D", Num: 5}
	d6 := Tile{Suit: "D", Num: 6}
	d7 := Tile{Suit: "D", Num: 7}
	d8 := Tile{Suit: "D", Num: 1}
	
	ts := []Tile{d7, d1, d6, d2, d5,d8}
	h := Hand{}
	h.New()
	h.Dots = append(h.Dots, ts...)
	slices.SortFunc(h.Dots, cmp)
	fmt.Printf("%v\n", h.Dots)

}
