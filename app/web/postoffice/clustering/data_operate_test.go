package clustering

import (
	"testing"

	"gameclustering.com/internal/core"
	"gameclustering.com/internal/persistence"
	"gameclustering.com/internal/protocol"
)

func TestDataOpt(t *testing.T) {
	core.CreateTestLog()
	local := persistence.BadgerLocal{InMemory: true, LogDisabled: true}
	err := local.Open()
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	defer local.Close()
	dsp := DataServiceProvider{Local: &local}
	k := []byte("key1")
	v := []byte("value1")
	data := protocol.Data{Key: k, Value: v, Header: &protocol.Header{FactoryId: 100, ClassId: 300}}
	dset := SetData{Data: &data, Opt: core.CREATE_DATA_REQUEST}
	err = dsp.create(dset)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	req := protocol.Request{Data: &data}
	dget := GetData{Request: &req}
	d, err := dsp.get(dget)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	if d.Header.Revision != 1 {
		t.Errorf("revision should be 1 %d", d.Header.Revision)
	}
	if d.Header.Size != 6 {
		t.Errorf("value size should be 6 %d", d.Header.Size)
	}
	if string(d.Value) != string(v) {
		t.Errorf("should be same values %s %s", string(d.Value), string(v))
	}
	err = dsp.create(dset)
	if err == nil {
		t.Error("should be an error")
	}
	core.AppLog.Debug().Msgf("ERROR :%s", err.Error())
	data.Header.Revision = 1
	data.Value = []byte("value123")
	dupdate := SetData{Data: &data, Opt: core.UPDATE_DATA_REQUEST}
	err = dsp.update(dupdate)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	req1 := protocol.Request{Data: &data}
	dget1 := GetData{Request: &req1}
	d1, err := dsp.get(dget1)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	if d1.Header.Revision != 2 {
		t.Errorf("revision should be 2 %d", d1.Header.Revision)
	}
	if d1.Header.Size != 8 {
		t.Errorf("value size should be 8 %d", d.Header.Size)
	}
	if string(d1.Value) != string("value123") {
		t.Errorf("should be same values %s %s", string(d1.Value), string(v))
	}
	data.Header.Revision = 100
	data.Value = []byte("value123678")
	rupdate := SetData{Data: &data, Opt: core.RESET_DATA_REQUEST}
	err = dsp.reset(rupdate)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}

	ddelete := SetData{Data: &data, Opt: core.DELETE_DATA_REQUEST}
	err = dsp.delete(ddelete)
	if err != nil {
		t.Errorf("should not be error %s", err.Error())
	}
	req2 := protocol.Request{Data: &data}
	dget2 := GetData{Request: &req2}
	_, err = dsp.get(dget2)
	if err == nil {
		t.Errorf("should be error")
	}
	core.AppLog.Debug().Msgf("DELETED %s", err.Error())
}
