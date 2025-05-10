package persistence

import (
	"testing"
)



func TestLocalStore(t *testing.T) {
	local :=LocalStore{InMemory: true,Path: "/home/yinghu/local"}
	err := local.Open()
	if err!= nil{
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	//local.Db.Update() 	
}