package event

type CloseEvent struct {
	EventObj `json:"-"`
}

func (s *CloseEvent) ClassId() int {
	return CLOSE_CID
}