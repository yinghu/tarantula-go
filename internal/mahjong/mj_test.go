package mahjong

import (
	"testing"
)

func TestSuit(t *testing.T) {
	d1 := Tile{}
	d1.From(DOTS1)
	if d1.Value != 1 {
		t.Errorf("value should be 1 %d", d1.Value)
	}
	if d1.Suit != "D" {
		t.Errorf("SUIT should be D %s", d1.Suit)
	}
	c9 := Tile{}
	c9.From(CHARACTER9)
	if c9.Value != 9 {
		t.Errorf("value should be 9 %d", c9.Value)
	}
	if c9.Suit != "C" {
		t.Errorf("SUIT should be C %s", c9.Suit)
	}

	f3 := Tile{}
	f3.From(F_CHRYSANTHEMUM)
	if f3.Value != 0 {
		t.Errorf("value should be 0 %d", f3.Value)
	}
	if f3.Suit != F_CHRYSANTHEMUM {
		t.Errorf("SUIT should be F %s", f3.Suit)
	}

}

func TestShuffle(t *testing.T) {
	s := Stack{}
	s.New()
	if len(s.Deck) != s.Size {
		t.Errorf("deck size error %d", len(s.Deck))
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
