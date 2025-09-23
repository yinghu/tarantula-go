package main

const (
	CLASSIC_MAX_FORMS int = 4
)

type ClassicMahjong struct {
	Deck
	East  Hand
	South Hand
	West  Hand
	North Hand
}

func (c *ClassicMahjong) New() {
	c.Deck = Deck{}
	c.Deck.New()
	c.East = Hand{}
	c.East.New()
	c.South = Hand{}
	c.South.New()
	c.West = Hand{}
	c.West.New()
	c.North = Hand{}
	c.North.New()
	c.Shuffle()
} 

