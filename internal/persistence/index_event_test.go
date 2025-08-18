package persistence

import (
	"fmt"
	"strings"
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
	for i := range 10 {
		id, _ := sfk.Id()
		ie := event.IndexEvent{Id: id, Tag: "test:"}
		ie.Index = []byte("hix")
		ie.Key = []byte(fmt.Sprintf("%s%d", "hix", i))
		err = local.Save(&ie)
		if err != nil {
			t.Errorf("save error %s", err.Error())
		}
	}
	px := BufferProxy{}
	px.NewProxy(100)
	px.WriteString(event.INDEX_ETAG)
	px.WriteString("test:")
	px.WriteInt32(3)
	px.Write([]byte("hix"))
	px.Flip()
	local.List(&px, func(k, v core.DataBuffer, rev uint64) bool {
		v.ReadInt32()
		v.ReadInt64()
		ix := event.IndexEvent{}
		ix.ReadKey(k)
		if string(ix.Index) != "hix" {
			t.Errorf("wrong index %s", string(ix.Index))
			return false
		}
		ix.Read(v)
		ik := string(ix.Key)
		if !strings.HasPrefix(ik, "hix") {
			t.Errorf("wrong index %s", ik)
			return false
		}
		return true
	})
}
