package protocol

type LoginObjectHandler func(e *LoginObject) error

type LoginObjectListener struct {
	Callback LoginObjectHandler
}

func (m *LoginObjectListener) OnTask(task *Task) error {
	var me LoginObject
	err := task.Transaction.Message.UnmarshalTo(&me)
	if err != nil {
		return err
	}
	return m.Callback(&me)
}
