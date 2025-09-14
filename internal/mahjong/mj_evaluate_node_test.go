package mahjong

import "testing"

func TestEvaluationNode(t *testing.T) {
	e := EvaluationNode{}
	e.New(CHOW)
	d1 := Tile{Suit: "D", Rank: 1}
	d2 := Tile{Suit: "D", Rank: 2}
	d3 := Tile{Suit: "D", Rank: 3}
	p := []Tile{d1, d2, d3}
	m := Meld{Tiles: p}
	e.Formed = append(e.Formed, m)
	occ := e.OccurrenceOfMeld()
	if occ[m.Name()]!=1{
		t.Errorf("occuurence of %s should be 1",m.Name())
	}
	e.Formed = append(e.Formed, m)
	occ = e.OccurrenceOfMeld()
	if occ[m.Name()]!=2{
		t.Errorf("occuurence of %s should be 2",m.Name())
	}
	
	e1 := EvaluationNode{}
	e1.New(CHOW)
	e1.Formed = append(e1.Formed, m)
	if e.FormedIdentically(e1){
		t.Errorf("should not be identical")
	}
	e1.Formed = append(e1.Formed, m)
	if !e.FormedIdentically(e1){
		t.Errorf("should be identical")
	} 
}
