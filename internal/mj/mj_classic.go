package mj

const (
	CLASSIC_MAX_FORMS int = 4	
)

type ClassicMahjong struct {
	Deck
	//East  Hand
	//South Hand
	//West  Hand
	//North Hand
	Evaluator
}

func (c *ClassicMahjong) New() {
	c.Deck = Deck{}
	c.Deck.New()
	//c.East = Hand{}
	//c.East.New()
	//c.South = Hand{}
	//c.South.New()
	//c.West = Hand{}
	//c.West.New()
	//c.North = Hand{}
	//c.North.New()
	c.Shuffle()

}

func (c *ClassicMahjong) Mahjong(h *Hand) bool {
	h.Formed = append(h.Formed, c.Evaluate(h)...)
	var eyeCount int
	var formed int
	for _, v := range h.Formed {
		if v.Eye() {
			eyeCount++
		}
		formed++
	}
	return eyeCount == 1 && formed == h.MaxClaims
}
