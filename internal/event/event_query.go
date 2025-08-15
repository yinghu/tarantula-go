package event

import "gameclustering.com/internal/core"

type Query struct {
	Tag   string     `json:"Tag"`
	Limit int32      `json:"Limit"`
	Cc    chan Chunk `json:"-"`
}

func (q *Query) Write(buff core.DataBuffer) error {
	buff.WriteString(q.Tag)
	return nil
}
