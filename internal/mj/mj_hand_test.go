package mj

import (
	"fmt"
	"testing"
)

func TestHandPung(t *testing.T) {
	h := Hand{}
	h.New()
	h.append(FromS("B1"), false)
	h.append(FromS("B2"), false)
	h.append(FromS("B3"), false)
	h.append(FromS("B3"), false)
	h.append(FromS("C2"), false)
	h.append(FromS("C3"), false)
	h.append(FromS("B3"), false)
	
	fmt.Printf("%v\n", h.Tiles)
	d := FromS("B3")
	err := h.Pung(d)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("After pung %v\n", h.Tiles)
	fmt.Printf("After pung %v\n", h.Formed)
	c := FromS("C1")
	c2 := FromS("C2")
	c3 := FromS("C3")
	err = h.Chow(c,c2,c3)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("After chow %v\n", h.Tiles)
	fmt.Printf("After chow %v\n", h.Formed) 
}
