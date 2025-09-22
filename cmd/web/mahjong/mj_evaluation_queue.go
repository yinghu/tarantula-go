package main

import "fmt"

type EvaluationQueue struct {
	PendingNode []EvaluationNode
	Formed      []Meld
}


func (q *EvaluationQueue) Poll() (EvaluationNode, error) {
	if len(q.PendingNode) == 0 {
		return EvaluationNode{}, fmt.Errorf("no node on queue")
	}
	e := q.PendingNode[0]
	q.PendingNode = q.PendingNode[1:]
	return e, nil
}

func (q *EvaluationQueue) Next(e EvaluationNode) {
	for _, n := range q.PendingNode {
		if n.FormedIdentically(e) {
			return
		}
	}
	q.PendingNode = append(q.PendingNode, e)
}
