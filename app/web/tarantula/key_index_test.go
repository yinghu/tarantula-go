package main

import (
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/tarantula/protocol"
	hash "github.com/spaolacci/murmur3"
)

func TestKeyIndex2(t *testing.T) {
	core.CreateTestLog()
	local := persistence.BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	dsp := DataServiceProvider{Local: &local}
	key := []byte("key1")
	kw := KeyIndex{Prefix: hash.Sum32(key), Key: key, Header: &protocol.Header{Revision: 1, FactoryId: 1, ClassId: 1, Timestamp: 1,Size: 100}}
	err = dsp.saveKeyIndex(&kw)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	kr := KeyIndex{Prefix: kw.Prefix, Key: kw.Key, Header: &protocol.Header{}}
	err = dsp.loadKeyIndex(&kr)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	if kw.Header.Revision != kr.Header.Revision {
		t.Errorf("should be same %d %d",kw.Header.Revision,kr.Header.Revision)
	}
	if kw.Header.FactoryId != kr.Header.FactoryId {
		t.Errorf("should be same %d %d",kw.Header.FactoryId,kr.Header.FactoryId)
	}
	if kw.Header.ClassId != kr.Header.ClassId {
		t.Errorf("should be same %d %d",kw.Header.ClassId,kr.Header.ClassId)
	}
	if kw.Header.Timestamp != kr.Header.Timestamp {
		t.Errorf("should be same %d %d",kw.Header.Timestamp,kr.Header.Timestamp)
	}
	if kw.Header.Size != kr.Header.Size {
		t.Errorf("should be same %d %d",kw.Header.Size,kr.Header.Size)
	}
}
