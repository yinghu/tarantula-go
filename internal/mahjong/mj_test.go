package mahjong

import (
	"testing"
)

func TestSuit(t *testing.T) {
	d1 := Tile{}
	d1.From(DOTS1)
	if d1.Num != 1 {
		t.Errorf("value should be 1 %d", d1.Num)
	}
	if d1.Suit != "D" {
		t.Errorf("SUIT should be D %s", d1.Suit)
	}
	c9 := Tile{}
	c9.From(CHARACTER9)
	if c9.Num != 9 {
		t.Errorf("value should be 9 %d", c9.Num)
	}
	if c9.Suit != "C" {
		t.Errorf("SUIT should be C %s", c9.Suit)
	}

	f3 := Tile{}
	f3.From(F_CHRYSANTHEMUM)
	if f3.Num != 0 {
		t.Errorf("value should be 0 %d", f3.Num)
	}
	if f3.Suit != F_CHRYSANTHEMUM {
		t.Errorf("SUIT should be F %s", f3.Suit)
	}

}

func TestShuffle(t *testing.T) {
	s := Deck{}
	s.New()
	if len(s.Stack) != s.Size {
		t.Errorf("deck size error %d", len(s.Stack))
	}
	s.Shuffle()
	if s.header != 0 {
		t.Errorf("deck header error %d", s.header)
	}
	if s.tail == s.Size {
		t.Errorf("deck tail error %d", s.tail)
	}
	_, err := s.Draw()
	if err != nil {
		t.Errorf("first draw erro %s", err.Error())
	}
	if s.header != 1 {
		t.Errorf("first draw erro %d", s.header)
	}
	tail := s.tail
	_, err = s.Kong()
	if err != nil {
		t.Errorf("first kong erro %s", err.Error())
	}
	if s.tail != tail-1 {
		t.Errorf("first kong erro %d", s.tail)
	}
	s.tail = 1
	_, err = s.Draw()
	if err != nil {
		t.Errorf("should be last draw %s", err.Error())
	}
	noTile, err := s.Draw()
	if err == nil {
		t.Errorf("should be no more draw %s", noTile.Suit)
	}
	s.header = 100
	s.tail = 100
	_, err = s.Kong()
	if err != nil {
		t.Errorf("should be last knog %s", err.Error())
	}
	noKnog, err := s.Kong()
	if err == nil {
		t.Errorf("should be no more know %s", noKnog.Suit)
	}
}

func TestEyeMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Num: 1}
	p := []Tile{d1, d1}
	m := Meld{Tiles: p}
	if !m.Eye() {
		t.Errorf("should be an eye")
	}
	m.Tiles = append(m.Tiles, d1)
	if m.Eye() {
		t.Errorf("should not be an eye")
	}
}

func TestChowMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Num: 1}
	d2 := Tile{Suit: "D", Num: 2}
	d3 := Tile{Suit: "D", Num: 3}
	p := []Tile{d1, d2, d3}
	m := Meld{Tiles: p}
	if !m.Chow() {
		t.Errorf("should be a chow")
	}
	m.Tiles = append(m.Tiles, d1)
	if m.Chow() {
		t.Errorf("should not be a chow")
	}
	x := []Tile{d1, d1, d2}
	c := Meld{Tiles: x}
	if c.Chow() {
		t.Errorf("should not be a chow")
	}
}

func TestPongMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Num: 1}
	p := []Tile{d1, d1, d1}
	m := Meld{Tiles: p}
	if !m.Pong() {
		t.Errorf("should be a pong")
	}
	m.Tiles = append(m.Tiles, d1)
	if m.Pong() {
		t.Errorf("should not be a pong")
	}
}

func TestKongMeld(t *testing.T) {
	d1 := Tile{Suit: "D", Num: 1}
	p := []Tile{d1, d1, d1, d1}
	m := Meld{Tiles: p}
	if !m.Kong() {
		t.Errorf("should be a kong")
	}
	m.Tiles = append(m.Tiles, d1)
	if m.Kong() {
		t.Errorf("should not be a kong")
	}
}

func TestHand(t *testing.T) {
	s := Deck{}
	s.New()
	s.Shuffle()
	h := Hand{}
	h.New()
	err := h.Draw(&s)
	if err != nil {
		t.Errorf("should be error")
	}
	err = h.Knog(&s)
	if err != nil {
		t.Errorf("should be error")
	}

}
