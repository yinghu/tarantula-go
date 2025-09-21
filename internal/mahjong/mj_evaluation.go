package mahjong

type Evaluation struct {
	Queue EvaluationQueue
}

func (e *Evaluation) Evaluate(h []Tile) []Meld {
	e.Queue = EvaluationQueue{PendingNode: make([]EvaluationNode, 0), Formed: make([]Meld, 0)}
	node := EvaluationNode{PendingHand: h, Formed: make([]Meld, 0), Type: NONE}
	e.Queue.Next(node)
	for {
		n, err := e.Queue.Poll()
		if err != nil {
			break
		}
		if len(n.PendingHand) == 0 && n.WellFormed() {
			e.Queue.Formed = append(e.Queue.Formed, n.Formed...)
			break
		}
		if len(n.Formed) == 4 {
			hi := HandIndex{}
			hi.From(n.PendingHand)
			_, err := hi.Eye()
			if err != nil {
				//e.Queue.Next()
			}
		}
		
	}
	return e.Queue.Formed
}
