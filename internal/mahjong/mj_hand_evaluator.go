package mahjong

type Evaluator struct {
	Queue EvaluationQueue
}

func (e *Evaluator) Evaluate(h *Hand) []Meld {
	e.Queue = EvaluationQueue{PendingNode: make([]EvaluationNode, 0), Formed: make([]Meld, 0)}
	node := EvaluationNode{PendingHand: h.Tiles, Formed: make([]Meld, 0), Type: NONE}
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
		hx := HandIndex{}
		hx.From(n.PendingHand)
		if len(n.Formed) == h.MaxForms {
			eye, err := hx.Eye()
			if err == nil {
				nx := EvaluationNode{PendingHand: hx.AfterFormed(eye), Formed: n.Form(eye), Type: EYE}
				e.Queue.Next(nx)
			}
		}

		kong := hx.Kong()
		for _, m := range kong {
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: KNOG}
			e.Queue.Next(nx)
		}

		pong := hx.Pong()
		for _, m := range pong {
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: PONG}
			e.Queue.Next(nx)
		}

		chow := hx.Chow()
		for _, m := range chow {
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: CHOW}
			e.Queue.Next(nx)
		}
		
	}
	return e.Queue.Formed
}
