package main

import "testing"

func TestMahjongEvent(t *testing.T) {
	err := NewMahjongErrorEvent(100, 200, 1, "error")
	if err.SystemId != 100 {
		t.Errorf("systemId should 100 %d", err.SystemId)
	}
	if err.RecipientId() != 100 {
		t.Errorf("recipientId should 100 %d", err.RecipientId())
	}
	if err.TableId != 200 {
		t.Errorf("table id should 200 %d", err.TableId)
	}
	if err.Code != 1 {
		t.Errorf("CODE should 1 %d", err.Code)
	}
	if err.Message != "error" {
		t.Errorf("error should error %s", err.Message)
	}
}
