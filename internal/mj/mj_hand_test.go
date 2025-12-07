package mj

import (
	"testing"
)

func TestHandPung(t *testing.T) {
	h := Hand{}
	h.New()
	h.append(FromS("B2"), false)
	h.append(FromS("B2"), false)
	h.append(FromS("B3"), false)
	h.append(FromS("B3"), false)
	h.append(FromS("C2"), false)
	h.append(FromS("C3"), false)
	h.append(FromS("B3"), false)

	d := FromS("B3")
	err := h.Pung(d)
	if err != nil {
		t.Errorf("%s\n", err.Error())
	}
	c := FromS("C1")
	c2 := FromS("C2")
	c3 := FromS("C3")
	err = h.Chow(c, c2, c3)
	if err != nil {
		t.Errorf("%s\n", err.Error())
	}
	err = h.Kong(d)
	if err != nil {
		t.Errorf("%s\n", err.Error())
	}
	//fmt.Printf("Tiles : %v\n", h.Tiles)
	//fmt.Printf("Melds : %v\n", h.Formed)
	h.append(FromS("D1"), false)
	h.append(FromS("D1"), false)
	h.append(FromS("D1"), false)
	err = h.Kong(FromS("D1"))
	if err != nil {
		t.Errorf("%s\n", err.Error())
	}
	//fmt.Printf("Tiles : %v\n", h.Tiles)
	//fmt.Printf("Melds : %v\n", h.Formed)
	h.append(FromS("D8"), false)
	h.append(FromS("D8"), false)
	h.append(FromS("D8"), false)
	h.append(FromS("D8"), false)
	err = h.Kong(FromS("D8"))
	if err != nil {
		t.Errorf("%s\n", err.Error())
	}
	//fmt.Printf("Tiles : %v\n", h.Tiles)
	//fmt.Printf("Melds : %v\n", h.Formed)
	mj := h.Mahjong()
	if !mj {
		t.Errorf("should be claimed %v", mj)
	}
}
