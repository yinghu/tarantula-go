package mj

import (
	"fmt"
	"testing"
)

type SampleListener struct {
}

func (s *SampleListener) OnIndex(t TileIndex) {
	fmt.Printf("tile usage %s : %d\n", t.Suit.Name(), t.Used)
}

func TestHandIndex(t *testing.T) {
	ix := HandIndex{Listener: &SampleListener{}}
	h := make([]Tile, 0)
	h = append(h, FromS("B1"))
	h = append(h, FromS("B2"))
	h = append(h, FromS("B3"))
	h = append(h, FromS("B2"))
	h = append(h, FromS("B3"))
	h = append(h, FromS("B4"))
	h = append(h, FromS("B3"))
	h = append(h, FromS("B3"))
	h = append(h, FromS("B6"))
	h = append(h, FromS("B6"))
	h = append(h, FromS("B6"))
	ix.From(h)
	ms := ix.Chow()
	for i := range ms {
		fmt.Printf("%s\n", ms[i].Name())
	}
	mp := ix.Pung()
	for i := range mp {
		fmt.Printf("%s\n", mp[i].Name())
	}

	mk := ix.Kong()
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
}
