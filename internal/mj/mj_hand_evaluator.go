package mj

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
		if len(n.Formed) == h.MaxClaims-len(h.Formed) {
			eye, err := hx.Eye()
			if err == nil {
				nx := EvaluationNode{PendingHand: hx.AfterFormed(eye), Formed: n.Form(eye), Type: EYE}
				e.Queue.Next(nx)
				if h.Listener != nil {
					h.Listener.OnFormed(eye)
				}
			}
		}

		kong := hx.Kong()
		for i := range kong {
			m := kong[i]
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: KNOG}
			e.Queue.Next(nx)
			if h.Listener != nil {
				h.Listener.OnFormed(m)
			}
		}

		pung := hx.Pung()
		for i := range pung {
			m := pung[i]
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: PUNG}
			e.Queue.Next(nx)
			if h.Listener != nil {
				h.Listener.OnFormed(m)
			}
		}

		chow := hx.Chow()
		for i := range chow {
			m := chow[i]
			nx := EvaluationNode{PendingHand: hx.AfterFormed(m), Formed: n.Form(m), Type: CHOW}
			e.Queue.Next(nx)
			if h.Listener != nil {
				h.Listener.OnFormed(m)
			}
		}

	}
	return e.Queue.Formed
}
