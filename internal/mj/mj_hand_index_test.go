package mj

import (
	"fmt"
	"testing"
)

type SampleListener struct {
	Usage map[int]int
}

func (s *SampleListener) OnIndex(t TileIndex) {
	u, e := s.Usage[t.Suit.Seq]
	if e {
		u = u + t.Used
		s.Usage[t.Suit.Seq] = u
	} else {
		s.Usage[t.Suit.Seq] = t.Used
	}
	//fmt.Printf("tile usage %s : %d\n", t.Suit.Name(), t.Used)
}

func TestHandIndex(t *testing.T) {
	sample := SampleListener{Usage: make(map[int]int)}
	ix := HandIndex{Listener: &sample}
	h := make([]Tile, 0)
	h = append(h, FromS("B1"))
	h = append(h, FromS("B2"))
	h = append(h, FromS("B3"))
	h = append(h, FromS("B2"))
	h = append(h, FromS("B5"))
	//h = append(h, FromS("B4"))
	//h = append(h, FromS("B3"))
	//h = append(h, FromS("B3"))
	//h = append(h, FromS("B6"))
	//h = append(h, FromS("B6"))
	h = append(h, FromS("B5"))
	ix.From(h)
	ix.Chow()
	ix.Pung()
	ix.Kong()
	ix.Eyes()
	for x := range sample.Usage {
		t := sample.Usage[x]
		fmt.Printf("Usage : %s %d\n", FromQ(x).Name(), t)
	}
	//ix.Eye()
	/**
	for i := range mk {
		fmt.Printf("%s\n", mk[i].Name())
	}
	mc := ix.CheckChow(FromS("B3"))
	for i := range mc {
		fmt.Printf("CH %s %d\n", mc[i].Name(), mc[i].Type())
	}
	mn := ix.CheckPung(FromS("B6"))

	fmt.Printf("PN %s %d\n", mn.Name(), mn.Type())

	mm := ix.CheckKong(FromS("B6"))

	fmt.Printf("PK %s %d\n", mm.Name(), mm.Type())

	m0 := ix.CheckPung(FromS("C1"))
	fmt.Printf("PK %s %d\n", m0.Name(), m0.Type())
	**/
}
