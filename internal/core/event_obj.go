package core

type EventObj struct {
	Callback EventListener `json:"-"`
	PersistentableObj
	ENodeId string `json:"nodeId"`
	ETag    string `json:"tag"`
	ETopic  string `json:"topic"`
	EOid    int64  `json:"oid,string"`
}

func (s *EventObj) OnNodeId(t string) {
	s.ENodeId = t
}

func (s *EventObj) NodeId() string {
	return s.ENodeId
}

func (s *EventObj) OnTag(t string) {
	s.ETag = t
}

func (s *EventObj) Tag() string {
	return s.ETag
}

func (s *EventObj) OnTopic(t string) {
	s.ETopic = t
}

func (s *EventObj) Topic() string {
	return s.ETopic
}

func (s *EventObj) Inbound(buff DataBuffer) error {
	return nil
}
func (s *EventObj) Outbound(buff DataBuffer) error {
	return nil
}
func (s *EventObj) OnListener(el EventListener) {
	s.Callback = el
}
func (s *EventObj) Listener() EventListener {
	return s.Callback
}

func (s *EventObj) OnOId(oid int64) {
	s.EOid = oid
}

func (s *EventObj) OId() int64 {
	return s.EOid
}

func (s *EventObj) RecipientId() int64 {
	return 0
}

func (s *EventObj) OnRecipientId(recipientId int64) {

}

func (s *EventObj) FactoryId() uint32 {
	return EVENT_FACTORY_ID
}
