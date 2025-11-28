package mj

import (
	"fmt"
	"testing"
)

func TestHandIndex(t *testing.T) {
	ix := HandIndex{}
	h := make([]Tile,0)
	h = append(h,FromS("B1"))
	h = append(h,FromS("B2"))
	h = append(h,FromS("B3"))
	h = append(h,FromS("B2"))
	h = append(h,FromS("B3"))
	h = append(h,FromS("B4"))
	h = append(h,FromS("B3"))
	h = append(h,FromS("B3"))
	ix.From(h)
	ms := ix.Chow()
	for i := range ms{
		fmt.Printf("%s\n",ms[i].Name())
	}
	mp := ix.Pung()
	for i := range mp{
		fmt.Printf("%s\n",mp[i].Name())
	}
	
	mk := ix.Kong()
	for i := range mk{
		fmt.Printf("%s\n",mk[i].Name())
	}

	//f := ix.AfterFormed()
}
