package protocol

type KeyValueHandler func(e *KeyValue) error

type KeyValueListener struct {
	Callback KeyValueHandler
}

func (m *KeyValueListener) OnTransaction(t *Transaction) error {
	return m.Callback(t.Object)
}
