package mj

import (
	"fmt"
	"testing"
)

func TestHand(t *testing.T) {
	cm := ClassicMahjong{}
	cm.New()
	h := Hand{}
	h.New()
	i := 14
	for{
		if i==0{
			break
		}
		h.Draw(&cm.Deck)
		i--
	}
	fmt.Printf("Hand %v\n",h.Tiles)
	//claimed := cm.Mahjong(&h)
	//if !claimed{
		//t.Errorf("should be claimed %v",claimed)
	//}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}
