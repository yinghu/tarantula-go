package persistence

import (
	"testing"

	buffer "github.com/0xc0d/encoding/bytebuffer"
)

func TestLocalStore(t *testing.T) {
	local := LocalStore{InMemory: true, Path: "/home/yinghu/local"}
	err := local.Open()
	if err != nil {
		t.Errorf("Local store error %s", err.Error())
	}
	defer local.Close()
	key := buffer.NewByteBuffer(10)
	key.PutUint16(10)
	key.Flip()
	value := buffer.NewByteBuffer(10)
	value.PutBytes([]byte("hello"),0,5)
	value.Flip()
	e := local.Set(*key, *value)
	if e != nil {
		t.Errorf("Local store set error %s", err.Error())
	}
	key.Rewind()
	local.Get(*key)
}
