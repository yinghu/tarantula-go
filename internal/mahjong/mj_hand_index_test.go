package mahjong

import (
	"fmt"
	"testing"
)

func TestHandIndex(t *testing.T) {

	b1 := NewTile(BAMBOO1)
	b2 := NewTile(BAMBOO2)
	b3 := NewTile(BAMBOO3)

	b4 := NewTile(BAMBOO2)
	b5 := NewTile(BAMBOO3)
	b6 := NewTile(BAMBOO4)

	b7 := NewTile(BAMBOO3)
	b8 := NewTile(BAMBOO4)
	b9 := NewTile(BAMBOO5)

	c7 := NewTile(BAMBOO4)
	c8 := NewTile(BAMBOO5)
	c9 := NewTile(BAMBOO6)

	p1 := NewTile(BAMBOO1)
	p2 := NewTile(BAMBOO1)

	tiles := []Tile{b1, b2, b3, b4, b5, b6, b7, b8, b9, c7, c8, c9, p1, p2}
	nd := EvaluationNode{PendingHand: tiles, Formed: make([]Meld, 0), Type: NONE}
	q := EvaluationQueue{PendingNode: make([]EvaluationNode, 0), Formed: make([]Meld, 0)}
	q.Next(nd)
	for {
		n, err := q.Poll()
		if err != nil {
			break
		}
		if len(n.PendingHand) == 0 && n.WellFormed() {
			q.Formed = append(q.Formed, n.Formed...)
			break
		}
		hx := HandIndex{}
		hx.From(n.PendingHand)
		if len(n.Formed) == 4 {
			fmt.Printf("Eval : %v\n", n.Formed)
			fmt.Printf("Eye : %v\n", n.PendingHand)
			eye, err := hx.Eye()
			if err == nil {
				nx := EvaluationNode{PendingHand: hx.AfterFormed(eye), Formed: n.Form(eye), Type: EYE}
				q.Next(nx)
			}
		}

		kong := hx.Kong()
		for _, m := range kong {
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: KNOG}
			q.Next(nx)
		}

		pong := hx.Pong()
		for _, m := range pong {
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: PONG}
			q.Next(nx)
		}

		chow := hx.Chow()
		for _, m := range chow {
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: CHOW}
			q.Next(nx)
		}

	}
	//for m := range q.Formed{
	fmt.Printf("Formed %v\n", q.Formed)
}
