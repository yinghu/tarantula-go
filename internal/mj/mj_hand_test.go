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
	h.append(NewTile("B1"),false)
	h.append(NewTile("B2"),false)
	h.append(NewTile("B3"),false)
	fmt.Printf("Hand %v\n",h.Tiles)
	tx := NewTile("B3")
	td,err := h.Discharge(tx.Seq)
	if err!=nil{
		fmt.Printf("Err %s\n",err.Error())
		return
	}
	fmt.Printf("Deleted %v %v\n",td,h.Tiles)
	
	//claimed := cm.Mahjong(&h)
	//if !claimed{
		//t.Errorf("should be claimed %v",claimed)
	//}
	//B1,B2,B3,B2,B3,B4,B3,B4,B5,B4,B5,B6,B1,B1
}
