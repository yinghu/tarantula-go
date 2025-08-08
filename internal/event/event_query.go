package event

type Query struct {
	Tag   string     `json:"Tag"`
	Limit int32      `json:"Limit"`
	Cc    chan Chunk `json:"-"`
}
