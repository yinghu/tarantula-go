package protocol

type Handler func(e *MessageEvent)

type MessageEventListener struct {
	Callback Handler
}

func (m *MessageEventListener) OnTopic(topic *Topic) error {
	var me MessageEvent
	err := topic.Event.Message.UnmarshalTo(&me)
	if err != nil {
		return err
	}
	m.Callback(&me)
	return nil
}
