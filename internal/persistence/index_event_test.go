package persistence

import (
	"fmt"
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/event"
	"gameclustering.com/internal/util"
)

func TestIndexEvent(t *testing.T) {
	local := BadgerLocal{InMemory: false, Path: "/home/yinghu/local/index"}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	sfk := util.NewSnowflake(1, util.EpochMillisecondsFromMidnight(2020, 1, 1))
	px := NewBuffer(100)

	for i := range 10 {
		id, _ := sfk.Id()
		ie := event.IndexEvent{Id: id, Tag: "test:"}
		px.Clear()
		px.WriteString("join:")
		px.WriteInt64(100)
		px.Flip()
		ie.WriteIndexKey(px)
		px.Clear()
		px.WriteInt64(200)
		px.Flip()
		ie.WriteIndexValue(px)
		err = local.Save(&ie)
		if err != nil {
			t.Errorf("save error %s %d", err.Error(), i)
		}
	}
	q := event.QIndex{IndexTag: "test:"}
	q.Tag = event.INDEX_ETAG
	px.Clear()
	px.WriteString("join:")
	px.WriteInt64(100)
	px.Flip()
	q.WriteIndexKey(px)
	px.Clear()
	q.QCriteria(px)
	px.Flip()
	local.List(px, func(k, v core.DataBuffer, rev uint64) bool {
		v.ReadInt32()
		v.ReadInt64()
		ix := event.IndexEvent{}
		ix.ReadKey(k)
		ix.Read(v)
		px.Clear()
		px.Write(ix.IndexKey)
		px.Flip()
		s, _ := px.ReadString()
		d, _ := px.ReadInt64()
		fmt.Printf("Index key %s %d %d\n", s, d, ix.Id)
		px.Clear()
		px.Write(ix.IndexValue)
		px.Flip()
		vd, _ := px.ReadInt64()
		fmt.Printf("Index value %d\n", vd)
		return true
	})
}
