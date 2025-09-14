package mahjong


type EvaluationNode struct {
	PendingHand []Tile
	Formed      []Meld
	Type        uint8
}

func (e *EvaluationNode) New(t uint8) {
	e.PendingHand = make([]Tile, 0)
	e.Formed = make([]Meld, 0)
	e.Type = t
}

func (e *EvaluationNode) FormedIdentically(n EvaluationNode) bool {
	occ1 := n.OccurrenceOfMeld()
	occ2 := e.OccurrenceOfMeld()
	for k, v := range occ2 {
		cv, exists := occ1[k]
		if !exists {
			return false
		}
		if cv != v {
			return false
		}
		delete(occ1, k)
	}
	return len(occ1) == 0
}

func (e *EvaluationNode) OccurrenceOfMeld() map[string]int {
	occ := make(map[string]int)
	for _, v := range e.Formed {
		o, exists := occ[v.Name()]
		if exists {
			occ[v.Name()] = o + 1
		} else {
			occ[v.Name()] = 1
		}
	}
	return occ
}
