package protocol

type MessageEventHandler func(e *MessageEvent) error

type MessageEventListener struct {
	Callback MessageEventHandler
}

func (m *MessageEventListener) OnTopic(topic *Topic) error {
	var me MessageEvent
	err := topic.Event.Message.UnmarshalTo(&me)
	if err != nil {
		return err
	}
	return m.Callback(&me)
}
