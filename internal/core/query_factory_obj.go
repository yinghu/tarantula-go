package core

type QueryFactoryObj struct {
	Q Query
}

func (f *QueryFactoryObj) Import(data []byte) (Query, error) {
	buff := NewBuffer(COMPOSIT_KEY_MAX)
	if err := buff.Write(data); err != nil {
		return f.Q, err
	}
	buff.Flip()
	if err := f.Q.QRead(buff); err != nil {
		return f.Q, err
	}
	return f.Q, nil
}
func (f *QueryFactoryObj) Export(query Query) ([]byte, error) {
	var v []byte
	buff := NewBuffer(COMPOSIT_KEY_MAX)
	if err := query.QWrite(buff); err != nil {
		return v, nil
	}
	buff.Flip()
	return buff.Read(0)
}

func (f *QueryFactoryObj) Query() Query {
	return f.Q
}

func (p *QueryFactoryObj) WriteKey(key DataBuffer) error {
	return nil
}

func (p *QueryFactoryObj) ReadKey(key DataBuffer) error {
	return nil
}
