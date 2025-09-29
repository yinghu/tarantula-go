package event

type KickoffEvent struct {
	Source   string
	SystemId int64
	EventObj
}

func (s *KickoffEvent) ClassId() int {
	return KICKOFF_CID
}

func (s *KickoffEvent) ETag() string {
	return KICKOFF_ETAG
}

func (s *KickoffEvent) RecipientId() int64 {
	return s.SystemId
}
