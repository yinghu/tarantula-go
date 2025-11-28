package mj

import (
	"fmt"
	"testing"
)

func TestHandPung(t *testing.T) {
	h := Hand{}
	h.New()
	h.append(NewTile("B1"), false)
	h.append(NewTile("B2"), false)
	h.append(NewTile("B3"), false)
	h.append(NewTile("B3"), false)
	h.append(NewTile("C2"), false)
	h.append(NewTile("C3"), false)
	h.append(NewTile("B3"), false)
	
	fmt.Printf("%v\n", h.Tiles)
	d := NewTile("B3")
	err := h.Pung(d)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("After pung %v\n", h.Tiles)
	fmt.Printf("After pung %v\n", h.Formed)
	c := NewTile("C1")
	c2 := NewTile("C2")
	c3 := NewTile("C3")
	err = h.Chow(c,c2,c3)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("After chow %v\n", h.Tiles)
	fmt.Printf("After chow %v\n", h.Formed) 
}
