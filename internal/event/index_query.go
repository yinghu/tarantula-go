package event

import "gameclustering.com/internal/core"

type QIndex struct {
	IndexTag     string
	IndexKey []byte
	QWithTag
}

func (q *QIndex) QCriteria(buff core.DataBuffer) error {
	buff.WriteString(q.Tag)
	buff.WriteString(q.IndexTag)
	buff.WriteInt32(int32(len(q.IndexKey)))
	buff.Write(q.IndexKey)
	return nil
}

func (s *QIndex) WriteIndexKey(buff core.DataBuffer) error {
	d, err := buff.Read(0)
	if err != nil {
		return err
	}
	s.IndexKey = d
	return nil
}
