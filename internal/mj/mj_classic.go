package mj

const (
	CLASSIC_MAX_FORMS int = 4
)

type ClassicMahjong struct {
	Deck
	Evaluator
}

func (c *ClassicMahjong) New() {
	c.Deck = Deck{}
	c.Deck.New()
	c.Shuffle()
	c.Queue = EvaluationQueue{PendingNode: make([]EvaluationNode, 0), Formed: make([]Meld, 0)}
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
	return eyeCount == 1 && formed == h.MaxClaims+1
}

func (c *ClassicMahjong) CheckMahjong(h *Hand, match Tile) bool {
	cm := Hand{}
	cm.New()
	cm.Tiles = append(cm.Tiles, h.Tiles...)
	cm.Tiles = append(cm.Tiles, match)
	return c.Mahjong(&cm)
}
